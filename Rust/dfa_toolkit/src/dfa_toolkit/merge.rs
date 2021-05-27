use crate::dfa_toolkit::state_partition::StatePartition;
use crate::dfa_toolkit::merge_data::MergeData;
use std::time::Instant;
use std::collections::HashSet;

// StatePairScore struct to store state pairs and their merge score.
pub struct StatePairScore {
    pub state_1: i32, // StateID for first state.
    pub state_2: i32, // StateID for second state.
    pub score: f32,   // Score of merge for given states.
}

// rpni_search deterministically merges all possible state pairs within red-blue framework.
// The first valid merge with respect to the rejecting examples is chosen.
// Returns the resultant state partition and search data when no more valid merges are possible.
// Used by the regular positive and negative inference (RPNI) algorithm.
pub fn rpni_search(original_state_partition: StatePartition) -> (StatePartition, MergeData){
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
    // Start timer.
    let start = Instant::now();

    // Vector to store red states.
    let mut red_states = vec![state_partition.starting_block()];
    // Generated vector of blue states from red states.
    let mut blue_states = update_blue_set(&mut state_partition, &red_states);

    // Iterate until blue set is empty.
    while blue_states.len() > 0 {
        // Remove and store first state in blue states queue.
        let blue_state = blue_states.pop().unwrap();

        // Set merged flag to false.
        let mut merged = false;

        // Iterate over every red state.
        for red_state in red_states.iter(){
            // Increment merge count.
            merge_data.attempted_merges_count += 1;

            // Check if states are mergeable.
            if copied_partition.merge_states(*red_state, blue_state) {
                // Increment valid merge count.
                merge_data.valid_merges_count += 1;

                // Copy changes to original state partition.
                state_partition.copy_changes_from(&copied_partition);

                // Set merged flag to true.
                merged = true;
                break
            }

            // Undo merges from copied partition.
            copied_partition.rollback_changes_from(&state_partition)
        }

        // If merged flag is false, add current blue state
        // to red states set and ordered set.
        if !merged {
            red_states.push(blue_state)
        }

        // Update red and blue states using update_red_set and update_blue_set function.
        update_red_set(&mut state_partition, &mut red_states);
        blue_states = update_blue_set(&mut state_partition, &red_states);
    }

    // Add duration to search data.
    merge_data.duration = start.elapsed();

    return (state_partition, merge_data);
}

pub fn update_red_set(state_partition: &mut StatePartition, red_states: &mut Vec<i32>) {
    // Step 1 - Gather root of old red states and store in map declared below.
    // Iterate over every red state.
    for red_state in red_states.iter_mut(){
        // Replace red state with its parent using Find.
        *red_state = state_partition.find(*red_state);
    }
}

// update_red_blue_sets updates the red and blue sets given the state partition and the red set within the Red-Blue framework
// such as the GeneralizedRedBlueMerging function. It returns the blue set and modifies the red set via its pointer. This is
// used when the state partition is changed or when new states have been added to the red set.
pub fn update_blue_set(state_partition: &mut StatePartition, red_states: &Vec<i32>) -> Vec<i32> {
    // Step 1 - Gather root of old red states and store in map declared below.

    // Initialize set of red root states (blocks) to empty set.
    let mut red_set = HashSet::with_capacity(red_states.len());

    // Iterate over every red state.
    for red_state in red_states{
        // Add red state to red set.
        red_set.insert(red_state);
    }

    // Step 2 - Gather all blue states and store in map declared below.

    // Initialize set of blue states to empty vector.
    let mut blue_states = vec![];

    // Iterate over every red state.
    for red_state in red_states{
        // Iterate over each symbol within DFA.
        for symbol in 0..state_partition.alphabet_size {
            // If transition is valid and resultant state is not red,
            // add the parent block of the resultant state to blue set.
            let mut resultant_state_id = state_partition.blocks[*red_state as usize].transitions[symbol];
            if resultant_state_id > -1 {
                resultant_state_id = state_partition.find(resultant_state_id);

                if !red_set.contains(&resultant_state_id) {
                    blue_states.push(resultant_state_id);
                }
            }
        }
    }

    // Return populated vectors of blue and red states.
    return blue_states
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

