use serde::Deserialize;

pub const REJECTING: i8 = 0;
pub const ACCEPTING: i8 = 1;
pub const UNLABELLED: i8 = 2;

pub type StateLabel = i8;

#[derive(Debug, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct State {
    pub label: StateLabel,
    pub transitions: Vec<i32>,
    // #[serde(skip)]
    // depth: i32,
    // #[serde(skip)]
    // order: i32
}
