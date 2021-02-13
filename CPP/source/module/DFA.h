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
};

class DFA {
public:
    map<unsigned int, State> states;
    State startingState;
    vector<char> alphabet;

    DFA(map<unsigned int, State>& states, State& startingState, vector<char>& alphabet);
    vector<State> getAcceptingStates();
    vector<State> getRejectingStates();
    void addState(StateStatus& stateStatus);
    unsigned int depth(); 
    void describe(bool detail);

private:
    void depthUtil(State& state, int count, map<unsigned int, unsigned int>& stateMap);
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
