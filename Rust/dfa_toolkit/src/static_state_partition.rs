use crate::dfa::DFA;
use crate::dfa_state::{StateLabel, ACCEPTING, REJECTING, UNKNOWN};

#[derive(Copy, Clone)]
pub struct StaticBlock {
    pub root: i32,
    pub size: i32,
    pub link: i32,
    pub label: StateLabel,
    pub changed: bool,
    pub transitions: [i32; 2],
}

pub struct StaticStatePartition {
    pub blocks: Vec<StaticBlock>,
    pub blocks_count: i32,
    pub accepting_blocks_count: i32,
    pub rejecting_blocks_count: i32,
    pub alphabet_size: usize,
    pub starting_state_id: i32,

    pub is_copy: bool,
    pub changed_blocks: Vec<i32>,
    pub changed_blocks_count: i32,
}

impl StaticStatePartition {
    pub fn new(dfa: &DFA) -> StaticStatePartition {
        let mut state_partition = StaticStatePartition {
            blocks: Vec::with_capacity(dfa.states.len()),
            blocks_count: dfa.states.len() as i32,
            accepting_blocks_count: 0,
            rejecting_blocks_count: 0,
            alphabet_size: dfa.symbols_count,
            starting_state_id: dfa.starting_state_id,
            is_copy: false,
            changed_blocks: vec![],
            changed_blocks_count: 0,
        };

        for i in 0..dfa.states.len() {
            let mut block = StaticBlock {
                root: i as i32,
                size: 1,
                link: i as i32,
                label: dfa.states[i].label,
                changed: false,
                transitions: [0, 0],
            };

            if block.label == ACCEPTING {
                state_partition.accepting_blocks_count += 1;
            } else if block.label == REJECTING {
                state_partition.rejecting_blocks_count += 1;
            }

            for j in 0..dfa.symbols_count {
                block.transitions[j] = dfa.states[i].transitions[j as usize];
            }

            state_partition.blocks.push(block);
        }

        return state_partition;
    }

    pub fn copy(&self) -> StaticStatePartition {
        return StaticStatePartition {
            blocks: self.blocks.clone(),
            blocks_count: self.blocks_count,
            accepting_blocks_count: self.accepting_blocks_count,
            rejecting_blocks_count: self.rejecting_blocks_count,
            alphabet_size: self.alphabet_size,
            starting_state_id: self.starting_state_id,
            is_copy: true,
            changed_blocks: vec![0; self.blocks.len()],
            changed_blocks_count: 0,
        };
    }

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
        if self.blocks[block_id_1 as usize].label == UNKNOWN
            && self.blocks[block_id_2 as usize].label != UNKNOWN
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

    pub fn rollback_changes(&mut self, original_state_partition: &StaticStatePartition) {
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
                self.blocks[block_id as usize] = original_state_partition.blocks[block_id as usize];
            }

            // Empty the changed blocks vector.
            self.changed_blocks_count = 0;
        }
    }
}

impl DFA {
    pub fn to_static_state_partition(&self) -> StaticStatePartition {
        return StaticStatePartition::new(self);
    }
}
