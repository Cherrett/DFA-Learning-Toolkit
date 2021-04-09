use crate::dfa::dfa_from_json;
use std::time::Instant;

mod dfa;
mod dfa_state;
mod state_partition;

fn main() {
    det_merge_benchmark();
}

fn det_merge_benchmark(){
    // These are target dfa sizes we will test.
    let dfa_sizes = [16, 32, 64];
    // These are the training set sizes we will test.
    let training_set_sizes = [230, 607, 1521];

    let mut iterator = 0;

    while iterator < dfa_sizes.len(){
        let target_size = dfa_sizes[iterator];
        let training_set_size = training_set_sizes[iterator];

        println!("-------------------------------------------------------------");
        println!("-------------------------------------------------------------");
        println!("BENCHMARK {} (Target: {} states, Training: {} strings", iterator+1, target_size, training_set_size);
        println!("-------------------------------------------------------------");
        println!("-------------------------------------------------------------");

        // Create APTA.
        let apta = dfa_from_json(format!("TestingAPTAs/{}.json", target_size));

        println!("APTA size: {}", apta.states.len());

        // Perform all the merges.
        let part = apta.to_state_partition();
        let mut snapshot = part.copy();
        let mut total_merges = 0;
        let mut valid_merges = 0;
        let start = Instant::now();

        let mut i: i32 = 0;
        while i < apta.states.len() as i32 {
            let mut j: i32 = i + 1;
            while j < apta.states.len() as i32 {
                total_merges += 1;
                if snapshot.merge_states(apta, i, j){
                    valid_merges += 1;
                }
                j += 1;
            }
            i += 1;
        }



        let total_time = start.elapsed().as_secs_f64();
        println!("Total merges: {}", total_merges);
        println!("Valid merges: {}", valid_merges);
        println!("Time: {}s", total_time);
        println!("Merges per second: {}", total_merges as f64/total_time);

        iterator += 1;
    }
}