use crate::dfa_toolkit::dfa::DFA;
use crate::dfa_toolkit::state_partition::StatePartition;
use crate::dfa_toolkit::merge::exhaustive_search_using_scoring_function;
use crate::dfa_toolkit::merge_data::MergeData;

// exhaustive_edsm is a greedy version of Evidence Driven State-Merging.
// It takes a DFA (APTA) as an argument which is used within the greedy search.
pub fn exhaustive_edsm(apta: DFA) -> (DFA, MergeData){
    // Store length of dataset.
    let length_of_dataset = apta.labelled_state_count();

    // EDSM scoring function.
    let edsm= |_state_id_1: i32, _state_id_2: i32, _partition_before: &StatePartition, partition_after: &StatePartition| -> f32{
        return (length_of_dataset - partition_after.number_of_labelled_blocks()) as f32;
    };

    // Convert APTA to StatePartition for state merging.
    let state_partition = apta.to_state_partition();

    // Call ExhaustiveSearchUsingScoringFunction function using state partition and EDSM scoring function
    // declared above. This function returns the resultant state partition and the search data.
    let mut result = exhaustive_search_using_scoring_function(state_partition, edsm);

    // Convert the state partition to a DFA.
    let resultant_dfa = result.0.to_quotient_dfa();

    // Check if DFA generated is valid.
    //resultant_dfa.is_valid_panic();

    // Return resultant DFA.
    return (resultant_dfa, result.1)
}