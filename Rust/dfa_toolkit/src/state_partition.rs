use crate::dfa_state::{StateLabel, UNKNOWN, ACCEPTING, REJECTING};
use std::collections::HashSet;
use crate::dfa::DFA;

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
    pub(crate) changed_blocks: HashSet<i32>
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
        changed_blocks: Default::default()
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
    pub fn copy(self) -> StatePartition{
        return StatePartition{
            blocks: self.blocks.clone(),
            blocks_count: self.blocks_count,
            accepting_blocks_count: self.accepting_blocks_count,
            rejecting_blocks_count: self.rejecting_blocks_count,
            is_copy: true,
            changed_blocks: Default::default()
        };
    }

    pub fn merge_states(&mut self, apta: DFA, state1: i32, state2: i32) -> bool{
        return true
    }
}

impl DFA{
    pub fn to_state_partition(&self) -> StatePartition{
        return new_state_partition(self);
    }
}