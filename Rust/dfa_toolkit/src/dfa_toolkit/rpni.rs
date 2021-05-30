use crate::dfa_toolkit::dfa::DFA;
use crate::dfa_toolkit::merge_data::MergeData;
use crate::dfa_toolkit::merge::rpni_search;

pub fn rpni(apta: DFA) -> (DFA, MergeData){
    // Convert APTA to StatePartition for state merging.
    let state_partition = apta.to_state_partition();

    // Call ExhaustiveSearchUsingScoringFunction function using state partition and EDSM scoring function
    // declared above. This function returns the resultant state partition and the search data.
    let (mut resultant_dfa, merge_data) = rpni_search(state_partition);

    // Convert the state partition to a DFA.
    let resultant_dfa = resultant_dfa.to_quotient_dfa();

    // Check if DFA generated is valid.
    resultant_dfa.is_valid_panic();

    // Return resultant DFA.
    return (resultant_dfa, merge_data)
}