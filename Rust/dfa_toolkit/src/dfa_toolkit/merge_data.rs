use crate::dfa_toolkit::merge::StatePairScore;
use std::time::Duration;

// MergeData struct to store merge search data.
pub struct MergeData{
    pub merges: Vec<StatePairScore>,
    pub attempted_merges_count: i32,
    pub valid_merges_count: i32,
    pub duration: Duration,
}

impl MergeData{
    pub fn attempted_merges_per_sec(&self) -> f64{
        return self.attempted_merges_count as f64/ self.duration.as_secs_f64()
    }
}