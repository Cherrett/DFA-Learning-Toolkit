use std::borrow::{Borrow, BorrowMut};
use crate::dfa_toolkit::dfa_state::{StateLabel, ACCEPTING, REJECTING, UNLABELLED};
use crate::dfa_toolkit::dfa::{DFA, new_dfa};
use std::collections::HashMap;

/// Block struct which represents a block within a partition.
pub struct Block {
    pub root: i32,              // Parent block of state.
    pub size: i32,              // Size (score) of block.
    pub link: i32,              // Index of next state within the block.
    pub label: StateLabel,      // Label of block.
    pub changed: bool,          // Whether block has been changed.
    pub transitions: Vec<i32>,  // Transition Table where each element corresponds to a transition for each symbol.
}

/// StatePartition struct which represents a State Partition.
pub struct StatePartition {
    pub blocks: Vec<Block>,             // Vector of blocks.
    pub blocks_count: i32,              // Number of blocks within partition.
    pub accepting_blocks_count: i32,    // Number of accepting blocks within partition.
    pub rejecting_blocks_count: i32,    // Number of rejecting blocks within partition.
    pub alphabet_size: usize,           // Size of alphabet.
    pub starting_state_id: i32,         // The ID of the starting state.

    pub is_copy: bool,                  // Whether state partition is a copy (for reverting merges).
    pub changed_blocks: Vec<i32>,       // Vector of changed blocks.
    pub changed_blocks_count: i32,      // Number of changed blocks within partition.
}

impl StatePartition {
    /// new returns an initialized State Partition.
    pub fn new(dfa: &DFA) -> StatePartition {
        let mut state_partition = StatePartition {
            blocks: Vec::with_capacity(dfa.states.len()),
            blocks_count: dfa.states.len() as i32,
            accepting_blocks_count: 0,
            rejecting_blocks_count: 0,
            alphabet_size: dfa.alphabet.len(),
            starting_state_id: dfa.starting_state_id,
            is_copy: false,
            changed_blocks: vec![],
            changed_blocks_count: 0,
        };

        for i in 0..dfa.states.len() {
            let mut block = Block {
                root: i as i32,
                size: 1,
                link: i as i32,
                label: dfa.states[i].label,
                changed: false,
                transitions: dfa.states[i].transitions.to_vec(),
            };

            if block.label == ACCEPTING {
                state_partition.accepting_blocks_count += 1;
            } else if block.label == REJECTING {
                state_partition.rejecting_blocks_count += 1;
            }

            block.transitions = dfa.states[i].transitions.to_vec();

            state_partition.blocks.push(block);
        }

        return state_partition;
    }

    /// changed_block updates the required fields to mark block as changed.
    pub fn changed_block(&mut self, block_id: i32) {
        // Check that state partition is a copy and that block is not modified.
        if self.is_copy && !self.blocks[block_id as usize].changed {
            // Update changed vector to include changed block ID.
            self.changed_blocks[self.changed_blocks_count as usize] = block_id;
            // Increment the changed blocks counter.
            self.changed_blocks_count += 1;
            // Set changed flag within block to true.
            self.blocks[block_id as usize].changed = true;
        }
    }

    /// union connects two blocks by comparing their respective
    /// size (score) values to keep the tree flat.
    pub fn union(&mut self, mut block_id_1: i32, mut block_id_2: i32) {
        // If state partition is a copy, call ChangedBlock for
        // both blocks so merge can be undone if necessary.
        self.changed_block(block_id_1);
        self.changed_block(block_id_2);

        // Decrement blocks count.
        self.blocks_count -= 1;

        // If size of parent node is smaller than size of child node, switch
        // parent and child nodes.
        if self.blocks[block_id_1 as usize].size < self.blocks[block_id_2 as usize].size {
            let temp_value = block_id_1;
            block_id_1 = block_id_2;
            block_id_2 = temp_value;
        }

        // Link nodes by assigning the link of parent to link of child and vice versa.
        let temp_value = self.blocks[block_id_1 as usize].link;
        self.blocks[block_id_1 as usize].link = self.blocks[block_id_2 as usize].link;
        self.blocks[block_id_2 as usize].link = temp_value;

        // Set root of child node to parent node.
        self.blocks[block_id_2 as usize].root = block_id_1;
        // Increment size (score) of parent node by size of child node.
        self.blocks[block_id_1 as usize].size += self.blocks[block_id_2 as usize].size;

        // If label of parent is unknown and label of child is
        // not unknown, set label of parent to label of child.
        if self.blocks[block_id_1 as usize].label == UNLABELLED
            && self.blocks[block_id_2 as usize].label != UNLABELLED
        {
            self.blocks[block_id_1 as usize].label = self.blocks[block_id_2 as usize].label;
        } else if self.blocks[block_id_1 as usize].label == ACCEPTING
            && self.blocks[block_id_2 as usize].label == ACCEPTING
        {
            // Else, if both blocks are accepting, decrement accepting blocks count.
            self.accepting_blocks_count -= 1;
        } else if self.blocks[block_id_1 as usize].label == REJECTING
            && self.blocks[block_id_2 as usize].label == REJECTING
        {
            // Else, if both blocks are rejecting, decrement rejecting blocks count.
            self.rejecting_blocks_count -= 1;
        }

        // Update transitions.
        for i in 0..self.alphabet_size {
            // If transition of parent does not exist,
            // set transition of child and set child
            // transition to -1 (remove transition).
            if self.blocks[block_id_1 as usize].transitions[i] == -1 {
                self.blocks[block_id_1 as usize].transitions[i] =
                    self.blocks[block_id_2 as usize].transitions[i];
                self.blocks[block_id_2 as usize].transitions[i] = -1;
            }
        }
    }

    /// find each root element while compressing the
    /// levels to find the root element of the stateID.
    pub fn find(&mut self, state: i32) -> i32 {
        let mut state_id = state;

        // Traverse each root block until state is reached.
        while self.blocks[state_id as usize].root != state_id {
            // Compress if necessary.
            if self.blocks[state_id as usize].root
                != self.blocks[self.blocks[state_id as usize].root as usize].root
            {
                // If compression is required, mark state as
                // changed and set root to root of parent.
                self.changed_block(state_id);
                self.blocks[state_id as usize].root =
                    self.blocks[self.blocks[state_id as usize].root as usize].root;
            }
            state_id = self.blocks[state_id as usize].root;
        }

        return state_id;
    }

    /// return_set returns the state IDs within given block.
    pub fn return_set(&self, block: i32) -> Vec<i32>{
        let mut block_id = block;

        // Hash set of state IDs and add root element to set.
        let mut block_elements: Vec<i32> = vec![block_id];

        // Add root element to set.
        //block_elements.push(block_id);
        // Set root to block ID.

        let root = block_id;

        // Iterate until link of current block ID is
        // not equal to the root block.
        while self.blocks[block_id as usize].link != root{
            // Set block ID to link of current block.
            block_id = self.blocks[block_id as usize].link;
            // Add block ID to block elements set.
            block_elements.push(block_id);
        }

        // Return state IDs within set.
        return block_elements;
    }

    /// within_same_block checks whether two states are within the same block.
    pub fn within_same_block(&mut self, state_id_1: i32, state_id_2: i32) -> bool {
        return self.find(state_id_1) == self.find(state_id_2)
    }

    /// to_quotient_dfa converts a State Partition to a quotient DFA and returns it.
    pub fn to_quotient_dfa(&mut self) -> DFA{
        // Hashmap to store corresponding new state for
        // each root block within state partition.
        let mut blocks_to_state_map = HashMap::new();

        // Initialize resultant DFA to be returned by function.
        let mut resultant_dfa = new_dfa();

        // Get root blocks within state partition.
        let root_blocks = self.root_blocks();

        // Create alphabet within DFA.
        for _ in 0..self.alphabet_size{
            resultant_dfa.add_symbol()
        }

        // Create a new state within DFA for each root block and
        // set state label to block label.
        for i in 0..root_blocks.len(){
            blocks_to_state_map.insert(root_blocks[i], resultant_dfa.add_state(self.blocks[root_blocks[i] as usize].label));
        }

        // Update transitions using transitions within blocks and block to state map.
        for i in 0..root_blocks.len(){
            for symbol in 0..self.alphabet_size{
                let resultant_state = self.blocks[root_blocks[i] as usize].transitions[symbol];
                if resultant_state > -1 {
                    resultant_dfa.states[*blocks_to_state_map.get(&root_blocks[i]).unwrap() as usize].transitions[symbol] = *blocks_to_state_map.get(&self.find(resultant_state)).unwrap();
                }
            }
        }

        // Set starting state using block to state map.
        resultant_dfa.starting_state_id = *blocks_to_state_map.get(&self.starting_block()).unwrap();

        // Return populated resultant DFA.
        return resultant_dfa
    }

    /// merge_states recursively merges states to merge state1 and state2. Returns false if merge
    /// results in a non-deterministic automaton. Returns true if merge was successful.
    pub fn merge_states(&mut self, state1: i32, state2: i32) -> bool {
        let mut state1_id = state1;
        let mut state2_id = state2;

        // If parent blocks (root) are the same as state ID, skip finding the root.
        // Else, find the parent block (root) using Find function.
        if self.blocks[state1_id as usize].root != state1_id {
            state1_id = self.find(state1_id);
        }
        if self.blocks[state2_id as usize].root != state2_id {
            state2_id = self.find(state2_id);
        }

        // Return true if states are already in the same block
        // since merge is not required.
        if state1_id == state2_id {
            return true;
        }

        // Get pointer of both blocks.
        //let (mut block1, mut block2) = (&self.blocks[state1_id as usize], &self.blocks[state2_id as usize]);

        // If labels are contradicting, return false since this results
        // in a non-deterministic automaton so merge cannot be done.
        if (self.blocks[state1_id as usize].label == ACCEPTING
            && self.blocks[state2_id as usize].label == REJECTING)
            || (self.blocks[state1_id as usize].label == REJECTING
                && self.blocks[state2_id as usize].label == ACCEPTING)
        {
            return false;
        }

        // Merge states within state partition.
        self.union(state1_id, state2_id);

        // Iterate over alphabet.
        for i in 0..self.alphabet_size {
            // If either block1 or block2 do not have a transition, continue
            // since no merge is required.
            if self.blocks[state1_id as usize].transitions[i] == -1
                || self.blocks[state2_id as usize].transitions[i] == -1
            {
                continue;
            }
            // Else, merge resultant blocks.
            if !self.merge_states(
                self.blocks[state1_id as usize].transitions[i],
                self.blocks[state2_id as usize].transitions[i],
            ) {
                return false;
            }
        }

        // Return true if this is reached (deterministic).
        return true;
    }

    /// clone returns a clone of the state partition.
    pub fn clone(&self) -> StatePartition{
        // Initialize new StatePartition struct using state partition.
        let mut cloned_state_partition = StatePartition {
            blocks: Vec::with_capacity(self.blocks.len()),
            blocks_count: self.blocks_count,
            accepting_blocks_count: self.accepting_blocks_count,
            rejecting_blocks_count: self.rejecting_blocks_count,
            alphabet_size: self.alphabet_size,
            starting_state_id: self.starting_state_id,
            is_copy: self.is_copy,
            changed_blocks: vec![],
            changed_blocks_count: self.changed_blocks_count,
        };

        // Copy blocks.
        for block in self.blocks.iter() {
            cloned_state_partition.blocks.push(Block {
                root: block.root,
                size: block.size,
                link: block.link,
                label: block.label,
                changed: block.changed,
                transitions: block.transitions.to_vec(),
            });
        }

        // If state partition is already a copy.
        if self.is_copy{
            // Copy changed blocks vector.
            cloned_state_partition.changed_blocks = self.changed_blocks.to_vec();
        }

        // Return cloned state partition.
        return cloned_state_partition;
    }

    /// copy returns a copy of the state partition in 'copy' (undo) mode.
    pub fn copy(&self) -> StatePartition {
        // Panic if state partition is already a copy.
        if self.is_copy{
            panic!("This state partition is already a copy.")
        }

        // Create a clone of the state partition.
        let mut copied_state_partition = self.clone();

        // Mark copied state partition as a copy.
        copied_state_partition.is_copy = true;
        copied_state_partition.changed_blocks_count = 0;
        copied_state_partition.changed_blocks = vec![0; self.blocks.len()];

        // Return copied state partition.
        return copied_state_partition;
    }

    ///rollback_changes_from reverts any changes made within state partition given the original state partition.
    pub fn rollback_changes_from(&mut self, original_state_partition: &StatePartition) {
        // If the state partition is a copy, copy values of changed blocks from original
        // state partition. Else, do nothing.
        if self.is_copy {
            // Set blocks count values to the original values.
            self.blocks_count = original_state_partition.blocks_count;
            self.accepting_blocks_count = original_state_partition.accepting_blocks_count;
            self.rejecting_blocks_count = original_state_partition.rejecting_blocks_count;

            // Iterate over each altered block (state).
            for i in 0..self.changed_blocks_count {
                let block_id = self.changed_blocks[i as usize];

                // Update root, size, link, and label.
                self.blocks[block_id as usize].root = original_state_partition.blocks[block_id as usize].root;
                self.blocks[block_id as usize].size = original_state_partition.blocks[block_id as usize].size;
                self.blocks[block_id as usize].link = original_state_partition.blocks[block_id as usize].link;
                self.blocks[block_id as usize].label = original_state_partition.blocks[block_id as usize].label;
                self.blocks[block_id as usize].changed = false;
                self.blocks[block_id as usize].transitions = original_state_partition.blocks[block_id as usize].transitions.to_vec();
            }

            // Empty the changed blocks vector.
            self.changed_blocks_count = 0;
        }
    }

    /// copy_changes_from copies the changes from one state partition to another and resets
    /// the changed values within the copied state partition.
    pub fn copy_changes_from(&mut self, copied_state_partition: &StatePartition){
        // If the state partition is a copy, copy values of changed blocks to original
        // state partition. Else, do nothing.
        if copied_state_partition.is_copy{
            // Set blocks count values to the new values.
            self.blocks_count = copied_state_partition.blocks_count;
            self.accepting_blocks_count = copied_state_partition.accepting_blocks_count;
            self.rejecting_blocks_count = copied_state_partition.rejecting_blocks_count;

            // Iterate over each altered block (state).
            for i in 0..copied_state_partition.changed_blocks_count {
                let block_id = copied_state_partition.changed_blocks[i as usize];
                // Get block pointer from copied state partition.
                let changed_block = copied_state_partition.blocks[block_id as usize].borrow();
                // Get block pointer from original state partition.
                let mut original_block = self.blocks[block_id as usize].borrow_mut();

                // Update root, size, link, and label.
                original_block.root = changed_block.root;
                original_block.size = changed_block.size;
                original_block.link = changed_block.link;
                original_block.label = changed_block.label;

                // Set changed property of changed block to false.
                original_block.changed = false;

                // Copy transitions from changed to original.
                original_block.transitions = changed_block.transitions.to_vec();
            }

            // Empty the changed blocks vector.
            self.changed_blocks_count = 0;
        }
    }

    /// number_of_labelled_blocks returns the number of labelled blocks (states) within state partition.
    pub fn number_of_labelled_blocks(&self) -> i32 {
        // Return the sum of the accepting and rejecting blocks count.
        return self.accepting_blocks_count + self.rejecting_blocks_count
    }

    /// root_blocks returns the IDs of root blocks as a vector of integers.
    pub fn root_blocks(&self) -> Vec<i32>{
        // Initialize vector using blocks count value.
        let mut root_blocks = vec![0; self.blocks_count as usize];
        // Index (count) of root blocks.
        let mut index = 0;

        // Iterate over each block within partition.
        for block_id in 0..self.blocks.len(){
            // Check if root of current block is equal to the block ID
            if self.blocks[block_id].root == block_id as i32{
                // Add to rootBlocks vector using index.
                root_blocks[index] = block_id as i32;
                // Increment index.
                index += 1;
                // If index is equal to blocks count,
                // break since all root blocks have
                // been found.
                if index == self.blocks_count as usize{
                    break
                }
            }
        }

        // Return populated vector.
        return root_blocks;
    }

    /// ordered_blocks returns the IDs of root blocks in order as a vector of integers
    pub fn ordered_blocks(&mut self) -> Vec<i32>{
        // Vector of boolean values to keep track of orders calculated.
        let mut order_computed = vec![false; self.blocks_count as usize];
        // Vector of integer values to keep track of ordered blocks.
        let mut ordered_blocks = vec![-1; self.blocks_count as usize];

        // Get starting block ID.
        let starting_block = self.starting_block();
        // Create a FIFO queue with starting block.
        let mut queue = vec![starting_block];
        // Add starting block to ordered blocks.
        ordered_blocks[0] = starting_block;
        // Mark starting block as computed.
        order_computed[starting_block as usize] = true;
        // Set index to 1.
        let mut index = 1;

        // Loop until queue is empty.
        while queue.len() > 0 {
            // Remove and store first state in queue.
            let block_id = queue.pop();

            // Iterate over each symbol (alphabet) within DFA.
            for symbol in 0..self.alphabet_size {
                // If transition from current state using current symbol is valid.
                let child_state_id = self.blocks[block_id as usize].transitions[symbol];

                if child_state_id >= 0 {
                    // Get block ID of child state.
                    let child_block_id = self.find(child_state_id);
                    // If depth for child block has been computed, skip block.
                    if !order_computed.contains(*child_block_id) {
                        // Add child block to queue.
                        queue.push(child_block_id);
                        // Set the order of the current block.
                        ordered_blocks[index] = child_block_id;
                        // Mark block as computed.
                        order_computed[child_block_id] = true;
                        // Increment current block order.
                        index += 1;
                    }
                }
            }
        }

        return ordered_blocks
    }

    /// starting_block returns the ID of the block which contains the starting state.
    pub fn starting_block(&mut self) -> i32{
        return self.find(self.starting_state_id)
    }
}

impl DFA {
    /// to_state_partition converts a DFA to a State Partition.
    pub fn to_state_partition(&self) -> StatePartition {
        return StatePartition::new(self);
    }
}
