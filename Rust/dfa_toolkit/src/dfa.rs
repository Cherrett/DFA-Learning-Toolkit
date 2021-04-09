use crate::dfa_state::State;
use std::fs;
use serde::Deserialize;

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct DFA{
    pub(crate) states: Vec<State>,
    pub(crate) starting_state_id: i32,
    pub(crate) symbols_count: i32,
    // #[serde(skip)]
    // depth: i32,
    // #[serde(skip)]
    // computed_depth_and_order: bool
}

pub fn dfa_from_json(file_path: String) -> DFA{
    let file = fs::File::open(file_path).expect("file should open read only");
    //let json: serde_json::Value = serde_json::from_reader(file).expect("file should be proper JSON");
    let dfa:DFA = serde_json::from_reader(file).expect("error while reading or parsing");

    return dfa;
}