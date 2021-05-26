use serde::Deserialize;
use std::fs;
use crate::dfa_toolkit::dfa_state::{State, ACCEPTING, REJECTING, StateLabel};

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct DFA {
    pub states: Vec<State>,
    pub starting_state_id: i32,
    pub alphabet: Vec<i32>,
    // #[serde(skip)]
    // depth: i32,
    // #[serde(skip)]
    // computed_depth_and_order: bool
}

impl DFA {
    pub fn add_state(&mut self, state_label: StateLabel) -> i32{
        // Create empty transition table with default values of -1 for each symbol within the DFA.
        let transitions = vec![-1 as i32, self.alphabet.len() as i32];
        // Initialize and add the new state to the slice of states within the DFA.
        self.states.push(State{ label: state_label, transitions });

        // Return the ID of the newly created state.
        return (self.states.len() - 1) as i32;
    }

    pub fn add_symbol(&mut self) {
        // Increment symbols count within the DFA.
        self.alphabet.push(self.alphabet.len() as i32);
        // Iterate over each state within the DFA and add an empty (-1) transition for the newly added state.
        for i in 0..self.states.len(){
            self.states[i].transitions.push(-1)
        }
    }

    // labelled_state_count returns the number of labelled states (accepting or rejecting) within DFA.
    pub fn labelled_state_count(&self) -> i32 {
        let mut count = 0;

        for state in self.states.iter(){
            if state.label == ACCEPTING || state.label == REJECTING{
                count += 1;
            }
        }

        return count;
    }
}

pub fn new_dfa() -> DFA {
    return DFA{
        states: vec![],
        starting_state_id: -1,
        alphabet: vec![]
    }
}

pub fn dfa_from_json(file_path: String) -> DFA {
    let file = fs::File::open(file_path).expect("file should open read only");
    //let json: serde_json::Value = serde_json::from_reader(file).expect("file should be proper JSON");
    let dfa: DFA = serde_json::from_reader(file).expect("error while reading or parsing");

    return dfa;
}
