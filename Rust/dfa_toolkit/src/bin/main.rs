extern crate dfa_toolkit;

use dfa_toolkit::dfa_toolkit::dfa::dfa_from_json;
use dfa_toolkit::dfa_toolkit::edsm::exhaustive_edsm;

fn main() {
    let apta = dfa_from_json(String::from("TestingAPTAs/32.json"));
    let result = exhaustive_edsm(apta);

    print!("Number of States: {}\n", result.0.states.len());
    print!("Duration: {:.2}s\n", result.1.duration.as_secs_f64());
    print!("Merges/s: {}\n", result.1.attempted_merges_per_sec().round());
    print!("Attempted Merges: {}\n", result.1.attempted_merges_count);
    print!("Valid Merges: {}\n", result.1.valid_merges_count);
}
