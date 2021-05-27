extern crate dfa_toolkit;

use dfa_toolkit::dfa_toolkit::dfa::dfa_from_json;
use dfa_toolkit::dfa_toolkit::rpni::rpni;

fn main() {
    let apta = dfa_from_json(String::from("TestingAPTAs/32.json"));
    let (dfa, merge_data) = rpni(apta);

    print!("Number of States: {}\n", dfa.states.len());
    print!("Duration: {:.2}s\n", merge_data.duration.as_secs_f64());
    print!("Merges/s: {}\n", merge_data.attempted_merges_per_sec().round());
    print!("Attempted Merges: {}\n", merge_data.attempted_merges_count);
    print!("Valid Merges: {}\n", merge_data.valid_merges_count);
}
