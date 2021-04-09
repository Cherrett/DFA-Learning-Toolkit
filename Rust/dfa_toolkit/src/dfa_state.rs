use serde::Deserialize;

pub const REJECTING:i8 = 0;
pub const ACCEPTING:i8 = 1;
pub const UNKNOWN:i8 = 2;

pub type StateLabel = i8;

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct State {
    pub(crate) label: StateLabel,
    pub(crate) transitions: Vec<i32>,
    // #[serde(skip)]
    // depth: i32,
    // #[serde(skip)]
    // order: i32
}
