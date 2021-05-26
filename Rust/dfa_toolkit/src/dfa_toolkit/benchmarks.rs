use std::time::Instant;
use crate::dfa_toolkit::dfa;

pub fn det_merge_benchmark() {
    println!("Normal State Partition");
    // These are target dfa sizes we will test.
    let dfa_sizes = [16, 32, 64];
    // These are the training set sizes we will test.
    let training_set_sizes = [230, 607, 1521];

    for iterator in 0..dfa_sizes.len() {
        let target_size = dfa_sizes[iterator];
        let training_set_size = training_set_sizes[iterator];

        println!("-------------------------------------------------------------");
        println!("-------------------------------------------------------------");
        println!(
            "BENCHMARK {} (Target: {} states, Training: {} strings",
            iterator + 1,
            target_size,
            training_set_size
        );
        println!("-------------------------------------------------------------");
        println!("-------------------------------------------------------------");

        // Create APTA.
        let apta = dfa::dfa_from_json(format!("TestingAPTAs/{}.json", target_size));

        println!("APTA size: {}", apta.states.len());

        // Perform all the merges.
        let part = apta.to_state_partition();
        let mut snapshot = part.copy();
        let mut total_merges = 0;
        let mut valid_merges = 0;
        let start = Instant::now();

        for i in 0..apta.states.len() as i32 {
            for j in (i + 1)..apta.states.len() as i32 {
                total_merges += 1;
                if snapshot.merge_states(i, j) {
                    valid_merges += 1;
                }

                snapshot.rollback_changes_from(&part);
            }
        }

        let total_time = start.elapsed().as_secs_f64();
        println!("Total merges: {}", total_merges);
        println!("Valid merges: {}", valid_merges);
        println!("Time: {:.2}s", total_time);
        println!("Merges per second: {:.2}\n", total_merges as f64 / total_time);
    }
}
