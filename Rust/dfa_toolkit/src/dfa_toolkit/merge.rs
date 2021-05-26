use crate::dfa_toolkit::state_partition::StatePartition;
use crate::dfa_toolkit::merge_data::MergeData;
use std::time::Instant;

// StatePairScore struct to store state pairs and their merge score.
pub struct StatePairScore {
    pub state_1: i32, // StateID for first state.
    pub state_2: i32, // StateID for second state.
    pub score: f32,   // Score of merge for given states.
}

// exhaustive_search_using_scoring_function deterministically merges all possible state pairs.
// The state pair to be merged is chosen using a scoring function passed as a parameter.
// Returns the resultant state partition and search data when no more valid merges are possible.
pub fn exhaustive_search_using_scoring_function(original_state_partition: StatePartition, scoring_function: impl Fn(i32, i32, &StatePartition, &StatePartition) -> f32) -> (StatePartition, MergeData){
    // Clone StatePartition.
    let mut state_partition = original_state_partition.clone();
    // Copy the state partition for undoing and copying changed states.
    let mut copied_partition = state_partition.copy();
    // Initialize search data.
    let mut merge_data = MergeData{
        merges: vec![],
        attempted_merges_count: 0,
        valid_merges_count: 0,
        duration: Default::default()
    };
    // State pair with the highest score.
    let mut highest_scoring_state_pair = StatePairScore{state_1: -1, state_2: -1, score: -1.0 };
    // Start timer.
    let start = Instant::now();

    // Loop until no more deterministic merges are available.
    loop {
        // Get root blocks within partition.
        let blocks = state_partition.root_blocks();

        // Get all valid merges and compute their score by
        // iterating over root blocks within partition.
        for i in 0..blocks.len(){
            for j in i+1..blocks.len(){
                // Increment merge count.
                merge_data.attempted_merges_count += 1;
                // Check if states are mergeable.
                if copied_partition.merge_states(blocks[i], blocks[j]){
                    // Increment valid merge count.
                    merge_data.valid_merges_count += 1;

                    // Calculate score.
                    let score = scoring_function(blocks[i], blocks[j], &state_partition, &copied_partition);

                    // If score is bigger than state pair with the highest score,
                    // set current state pair to state pair with the highest score.
                    if score > highest_scoring_state_pair.score{
                        highest_scoring_state_pair = StatePairScore{
                            state_1: blocks[i],
                            state_2: blocks[j],
                            score
                        }
                    }
                }

                // Undo merges from copied partition.
                copied_partition.rollback_changes_from(&state_partition)
            }
        }

        // Check if any deterministic merges were found.
        if highest_scoring_state_pair.score >= 0.0{
            // Merge the state pairs with the highest score.
            copied_partition.merge_states(highest_scoring_state_pair.state_1, highest_scoring_state_pair.state_2);
            // Copy changes to original state partition.
            state_partition.copy_changes_from(&copied_partition);

            // Add merged state pair with score to search data.
            merge_data.merges.push(highest_scoring_state_pair);

            // Remove previous state pair with the highest score.
            highest_scoring_state_pair = StatePairScore{state_1: -1, state_2: -1, score: -1.0 };
        }else{
            break;
        }
    }

    // Add duration to search data.
    merge_data.duration = start.elapsed();

    return (state_partition, merge_data);
}

