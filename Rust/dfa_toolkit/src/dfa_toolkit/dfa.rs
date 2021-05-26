use serde::Deserialize;
use std::fs;
use crate::dfa_toolkit::dfa_state::{State, ACCEPTING, REJECTING, StateLabel};
use std::collections::{HashSet};

#[derive(Debug, Deserialize)]
#[serde(rename_all = "PascalCase")]
/// DFA struct which represents a DFA.
pub struct DFA {
    pub states: Vec<State>,     // Vector of states within the DFA where the index is the State ID.
    pub starting_state_id: i32, // Alphabet within DFA.
    pub alphabet: Vec<i32>,     // The ID of the starting state of the DFA.
    // #[serde(skip)]
    // depth: i32,
    // #[serde(skip)]
    // computed_depth_and_order: bool
}

impl DFA {
    /// add_state adds a new state to the DFA with the corresponding State Label.
    /// Returns the new state's ID (index).
    pub fn add_state(&mut self, state_label: StateLabel) -> i32{
        // Create empty transition table with default values of -1 for each symbol within the DFA.
        let transitions = vec![-1 as i32, self.alphabet.len() as i32];
        // Initialize and add the new state to the vector of states within the DFA.
        self.states.push(State{ label: state_label, transitions });

        // Return the ID of the newly created state.
        return (self.states.len() - 1) as i32;
    }

    /// add_symbol adds a new symbol to the DFA.
    pub fn add_symbol(&mut self) {
        // Increment symbols count within the DFA.
        self.alphabet.push(self.alphabet.len() as i32);
        // Iterate over each state within the DFA and add an empty (-1) transition for the newly added state.
        for i in 0..self.states.len(){
            self.states[i].transitions.push(-1)
        }
    }

    /// labelled_state_count returns the number of labelled states (accepting or rejecting) within DFA.
    pub fn labelled_state_count(&self) -> i32 {
        let mut count = 0;

        for state in self.states.iter(){
            if state.label == ACCEPTING || state.label == REJECTING{
                count += 1;
            }
        }

        return count;
    }

    /// unreachable_states returns the state IDs of unreachable states. Extracted from:
    /// P. Linz, An Introduction to Formal Languages and Automata. Jones & Bartlett Publishers, 2011.
    pub fn unreachable_states(&self) -> Vec<i32> {
        // Hash set of reachable states made up of starting state.
        let mut reachable_states = HashSet::new();
        reachable_states.insert(self.starting_state_id);
        // Hash set of current states made up of starting state.
        let mut current_states = HashSet::new();
        current_states.insert(self.starting_state_id);

        // Iterate until current states is empty.
        while current_states.len() > 0 {
            // Hash set of next states.
            let mut next_states = HashSet::new();
            // Iterate over current states.
            for state_id in &current_states{
                // Iterate over each symbol within DFA.
                for symbol in 0..self.alphabet.len(){
                    // If transition from current state using current symbol
                    // is valid, add resultant state to next states.
                    let resultant_state_id = self.states[*state_id as usize].transitions[symbol as usize];
                    if resultant_state_id >= 0 {
                        next_states.insert(resultant_state_id);
                    }
                }
            }

            // Remove all state IDs from current states.
            current_states.clear();

            // Iterate over next states.
            for state_id in &next_states{
                // If state is not in reachable states map, add to
                // current states and to reachable states.
                // Else, ignore since state is already reachable.
                if !reachable_states.contains(state_id) {
                    current_states.insert(*state_id);
                    reachable_states.insert(*state_id);
                }
            }
        }

        // Vector of unreachable states.
        let mut unreachable_states = Vec::new();

        // Iterate over each state within DFA.
        for state_id in 0..self.states.len() {
            // If state ID is not in reachable states map,
            // add to unreachable states vector.
            if !reachable_states.contains(&(state_id as i32)) {
                unreachable_states.push(state_id as i32);
            }
        }

        return unreachable_states;
    }

    /// is_valid_panic checks whether DFA is valid.
    /// Panics if not valid. Used for error checking.
    pub fn is_valid_panic(&self) {
        if self.states.len() < 1 {
            // Panic if number of states is invalid.
            panic!("DFA does not contain any states.")
        }else if self.starting_state_id < 0 || self.starting_state_id >= self.states.len() as i32 {
            // Panic if starting state is invalid.
            panic!("Invalid starting state.")
        }else if self.alphabet.len() < 1 {
            // Panic if number of symbols is invalid.
            panic!("DFA does not contain any symbols.")
        }else if self.unreachable_states().len() > 0 {
            // Panic if any unreachable states exist within DFA.
            panic!("Unreachable State exist within DFA.")
        }
    }
}

/// new_dfa initializes a new empty DFA.
pub fn new_dfa() -> DFA {
    return DFA{
        states: vec![],
        starting_state_id: -1,
        alphabet: vec![]
    }
}

/// dfa_from_json returns a DFA read from a JSON file given a file path.
pub fn dfa_from_json(file_path: String) -> DFA {
    let file = fs::File::open(file_path).expect("file should open read only");
    let dfa: DFA = serde_json::from_reader(file).expect("error while reading or parsing");

    return dfa;
}
