#pragma once
#include <iostream>
#include <vector>
#include <algorithm>
#include <map>
using std::string;
using std::vector;
using std::map;
using std::sort;

enum class StateStatus {
    ACCEPTING = 1,
    REJECTING = 0,
    UNKNOWN = 2
};

class State {
public:
    StateStatus stateStatus;
    unsigned int stateID;
    map<char, unsigned int> transitions;

    State() = default;
    State(StateStatus stateStatus, unsigned int stateID);
    State(StateStatus stateStatus, unsigned int stateID, map<char, unsigned int> transitions);
};

class NFA_State {
public:
    StateStatus stateStatus;
    unsigned int stateID;
    map<char, vector<unsigned int>> transitions;

    NFA_State() = default;
    NFA_State(StateStatus stateStatus, unsigned int stateID);
    NFA_State(StateStatus stateStatus, unsigned int stateID, map<char, vector<unsigned int>> transitions);
};

//class TransitionFunction {
//public:
//    State fromState, toState;
//    char symbol;
//
//    TransitionFunction(State& fromState, State& toState, char& symbol);
//};

class DFA {
public:
    map<unsigned int, State> states;
    //vector<State> states;
    State startingState;
    vector<char> alphabet;
    //vector<TransitionFunction> transitionFunctions;

    DFA(map<unsigned int, State>& states, State& startingState, vector<char>& alphabet);
    vector<State> getAcceptingStates();
    vector<State> getRejectingStates();
    void addState(StateStatus& stateStatus);
    //void addState(State& state);
    //void addTransitionFunction(State& fromState, State& toState, char& symbol);
    unsigned int depth();
    void describe(bool detail);

private:
    void depthUtil(State& state, int count, map<unsigned int, unsigned int>& stateMap);
};

class NFA {
public:
    map<unsigned int, NFA_State> states;
    NFA_State startingState;
    vector<char> alphabet;

    NFA(map<unsigned int, NFA_State>& states, NFA_State& startingState, vector<char>& alphabet);
};

class StringInstance {
public:
    string stringValue;
    StateStatus stringStatus;
    unsigned int length;

    StringInstance(string& stringValue, StateStatus stringStatus, unsigned int& length);
    StringInstance(string& text, const string& delimiter);
    bool operator< (const StringInstance& otherString) const;
};

vector<StringInstance> GetListOfStringInstancesFromFile(string fileName);

void SortListOfStringInstancesInternal(vector<StringInstance>& strings);
vector<StringInstance> SortListOfStringInstances(vector<StringInstance> strings);

// Get (Augmented) Prefix Tree Acceptor from list of Strings
DFA GetPTAFromListOfStringInstances(vector<StringInstance>& strings, bool APTA);

bool StringInstanceConsistentWithDFA(StringInstance& string, DFA& dfa);

bool ListOfStringInstancesConsistentWithDFA(vector<StringInstance>& strings, DFA& dfa);

StateStatus GetStringStatusInRegardToDFA(StringInstance& string, DFA& dfa);

vector<StringInstance> GetAcceptingStringInstances(vector<StringInstance> strings);

vector<StringInstance> GetRejectingStringInstances(vector<StringInstance> strings);

bool isNFAdeterministic(NFA nfa);

NFA RPNI_Derive(DFA dfa, vector<vector<unsigned int>> partition);

struct RPNI_Deterministic_Merge_object {
public:
    vector<vector<unsigned int>> partition;
    DFA dfa;

    RPNI_Deterministic_Merge_object(vector<vector<unsigned int>> partition, DFA dfa);
};

NFA RPNI_Merge(NFA nfa, unsigned int state1, unsigned int state2);

RPNI_Deterministic_Merge_object RPNI_Deterministic_Merge(NFA nfa, vector<vector<unsigned int>> partition);

bool RPNI_StringInstanceConsistentWithDFA(StringInstance& string, DFA& dfa);

bool RPNI_ListOfNegativeStringInstancesConsistentWithDFA(vector<StringInstance>& strings, DFA& dfa);

DFA NFAtoDFA(NFA nfa);

DFA RPNI(vector<StringInstance>& acceptingStrings, vector<StringInstance>& rejectingStrings);