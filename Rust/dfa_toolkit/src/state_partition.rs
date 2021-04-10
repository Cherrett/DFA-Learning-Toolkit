use crate::dfa::DFA;
use crate::dfa_state::{StateLabel, UNKNOWN, ACCEPTING, REJECTING};

#[derive(Copy, Clone)]
pub struct Block{
    pub(crate) root: i32,
    pub(crate) size: i32,
    pub(crate) link: i32,
    pub(crate) label: StateLabel,
    pub(crate) changed: bool
}

pub struct StatePartition{
    pub(crate) blocks: Vec<Block>,
    pub(crate) blocks_count: i32,
    pub(crate) accepting_blocks_count: i32,
    pub(crate) rejecting_blocks_count: i32,

    pub(crate) is_copy: bool,
    pub(crate) changed_blocks: Vec<i32>,
    pub(crate) changed_blocks_count: i32
}

pub fn new_state_partition(dfa: &DFA) -> StatePartition{
    let mut state_partition = StatePartition{
        blocks: vec![Block{
            root: 0,
            size: 1,
            link: 0,
            label: UNKNOWN,
            changed: false
        }; dfa.states.len()],
        blocks_count: 0,
        accepting_blocks_count: 0,
        rejecting_blocks_count: 0,
        is_copy: false,
        changed_blocks: vec![],
        changed_blocks_count: 0
    };

    let mut i = 0;

    while i < dfa.states.len() {
        state_partition.blocks[i].root = i as i32;
        state_partition.blocks[i].link = i as i32;
        state_partition.blocks[i].label = dfa.states.get(i).unwrap().label;

        if state_partition.blocks[i].label == ACCEPTING{
            state_partition.accepting_blocks_count += 1;
        }else if state_partition.blocks[i].label == REJECTING{
            state_partition.rejecting_blocks_count += 1;
        }

        i += 1;
    }

    return state_partition
}

impl StatePartition{
    pub fn copy(&self) -> StatePartition{
        return StatePartition{
            blocks: self.blocks.clone(),
            blocks_count: self.blocks_count,
            accepting_blocks_count: self.accepting_blocks_count,
            rejecting_blocks_count: self.rejecting_blocks_count,
            is_copy: true,
            changed_blocks: vec![0;self.blocks.len()],
            changed_blocks_count: 0
        };
    }

    pub fn changed_block(&mut self, block_id: i32){
        // If block is not already changed.
        if !self.blocks[block_id as usize].changed{
            // Update changed vector to include changed block ID.
            self.changed_blocks[self.changed_blocks_count as usize] = block_id;
            // Increment the changed blocks counter.
            self.changed_blocks_count += 1;
            // Set changed flag within block to true.
            self.blocks[block_id as usize].changed = true;
        }
    }

    pub fn union(&mut self, block_id_1: i32, block_id_2: i32){
        // If state partition is a copy, call ChangedBlock for
        // both blocks so merge can be undone if necessary.
        if self.is_copy{
            self.changed_block(block_id_1);
            self.changed_block(block_id_2);
        }

        // Decrement blocks count.
        self.blocks_count -= 1;

        // Set block 1 to parent and block 2 to child.
        let (mut parent, mut child) = (block_id_1, block_id_2);

        // If size of parent node is smaller than size of child node, switch
        // parent and child nodes.
        if self.blocks[parent as usize].size < self.blocks[child as usize].size{
            let values = (child, parent);
            parent = values.0;
            child = values.1;
        }

        // Link nodes by assigning the link of parent to link of child and vice versa.
        let values = (self.blocks[child as usize].link, self.blocks[parent as usize].link);
        self.blocks[parent as usize].link = values.0;
        self.blocks[child as usize].link = values.1;

        // Get label of each block.
        let parent_label = self.blocks[parent as usize].label;
        let child_label = self.blocks[child as usize].label;

        // Set root of child node to parent node.
        self.blocks[child as usize].root = parent;
        // Increment size (score) of parent node by size of child node.
        self.blocks[parent as usize].size += self.blocks[child as usize].size;

        // If label of parent is unknown and label of child is
        // not unknown, set label of parent to label of child.
        if parent_label == UNKNOWN && child_label != UNKNOWN {
            self.blocks[parent as usize].label = child_label;
        } else if parent_label == ACCEPTING && child_label == ACCEPTING {
            // Else, if both blocks are accepting, decrement accepting blocks count.
            self.accepting_blocks_count -= 1;
        } else if parent_label == REJECTING && child_label == REJECTING {
            // Else, if both blocks are rejecting, decrement rejecting blocks count.
            self.rejecting_blocks_count -= 1;
        }
    }

    pub fn find(&mut self, state: i32) -> i32{
        let mut state_id = state;

        // Traverse each root block until state is reached.
        while self.blocks[state_id as usize].root != state_id {
            // Compress root.
            self.blocks[state_id as usize].root = self.blocks[self.blocks[state_id as usize].root as usize].root;
            state_id = self.blocks[state_id as usize].root;
        }

        return state_id;
    }

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

    pub fn merge_states(&mut self, dfa: &DFA, state1: i32, state2: i32) -> bool{
        let mut state1_id = state1;
        let mut state2_id = state2;

        // If parent blocks (root) are the same as state ID, skip finding the root.
        // Else, find the parent block (root) using Find.
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

        // Get label of each block.
        let state1_label = self.blocks[state1_id as usize].label;
        let state2_label = self.blocks[state2_id as usize].label;

        // If labels are contradicting, return false since this results
        // in a non-deterministic automaton so merge cannot be done.
        if (state1_label == ACCEPTING && state2_label == REJECTING) || (state1_label == REJECTING && state2_label == ACCEPTING) {
            return false;
        }

        let block1_set = self.return_set(state1_id);
        let block2_set = self.return_set(state2_id);

        // Merge states within state partition.
        self.union(state1_id, state2_id);

        // Iterate over each symbol within DFA.
        let mut symbol_id: i32 = 0;
        while symbol_id < dfa.symbols_count {
            // Iterate over each state within first block.
            for state_id in &block1_set{
                // Store resultant state from state transition of current state.
                let mut current_resultant_state_id = dfa.states[*state_id as usize].transitions[symbol_id as usize];

                // If resultant state ID is bigger than -1 (valid transition), get
                // the block containing state and store in transitionResult. The
                // states within the second block are then iterated and checked
                // for non-deterministic properties.
                if current_resultant_state_id > -1{
                    // Set resultant state to state transition for current symbol.
                    let transition_result = current_resultant_state_id;

                    // Iterate over each state within second block.
                    for state_id2 in &block2_set{
                        // Store resultant state from state transition of current state.
                        current_resultant_state_id = dfa.states[*state_id2 as usize].transitions[symbol_id as usize];
                        // If resultant state ID is bigger than -1 (valid transition), get the
                        // block containing state and compare it to the transition found above.
                        // If they are not equal, merge blocks to eliminate non-determinism.
                        if current_resultant_state_id > -1 {
                            // If resultant block is not equal to the block containing the state within transition
                            // found above, merge the two states to eliminate non-determinism.
                            // Merge states and if states cannot be merged, return false.
                            if !self.merge_states(dfa, transition_result, current_resultant_state_id){
                                return false;
                            }
                            // The loop is broken since the transition for the current symbol was found.
                            break;
                        }
                    }
                    // The loop is broken since the transition for the current symbol was found.
                    break
                }
            }

            symbol_id += 1;
        }

        // Return true if this is reached (deterministic).
        return true;
    }

    pub fn rollback_changes(&mut self, original_state_partition: &StatePartition){
        // If the state partition is a copy, copy values of changed blocks from original
        // state partition. Else, do nothing.
        if self.is_copy{
            // Set blocks count values to the original values.
            self.blocks_count = original_state_partition.blocks_count;
            self.accepting_blocks_count = original_state_partition.accepting_blocks_count;
            self.rejecting_blocks_count = original_state_partition.rejecting_blocks_count;

            // Iterate over each altered block (state).
            let mut i = 0;

            while i < self.changed_blocks_count{
                let state_id = self.changed_blocks[i as usize];
                self.blocks[state_id as usize] = original_state_partition.blocks[state_id as usize];

                i += 1;
            }

            // Empty the changed blocks vector.
            self.changed_blocks_count = 0;
        }
    }
}

impl DFA{
    pub fn to_state_partition(&self) -> StatePartition{
        return new_state_partition(self);
    }
}