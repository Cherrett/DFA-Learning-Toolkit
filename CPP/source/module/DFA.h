#pragma once
#include <iostream>
#include <vector>
#include <algorithm>
using std::string;
using std::vector;
using std::sort;

enum class StateStatus;

class State {
public:
    StateStatus stateStatus;
    unsigned int stateID;

    State() = default;
    State(StateStatus stateStatus, unsigned int stateID);
};

class TransitionFunction {
public:
    State fromState, toState;
    char symbol;

    TransitionFunction(State& fromState, State& toState, char& symbol);
};

class DFA {
public:
    vector<State> states;
    State startingState;
    vector<char> alphabet;
    vector<TransitionFunction> transitionFunctions;

    DFA(vector<State>& states, State& startingState, vector<char>& alphabet, vector<TransitionFunction>& transitionFunctions);

    vector<State> getAcceptingStates();

    void addState(StateStatus& stateStatus, unsigned int& statusID);

    void describe(bool detail);
};

class StringInstance {
public:
    string stringValue;
    bool accepting;
    unsigned int length;

    StringInstance(string& stringValue, bool accepting, unsigned int& length);
    StringInstance(string& text, const string& delimiter);
    bool operator< (const StringInstance& otherString) const;
};

vector<StringInstance> GetListOfStringInstancesFromFile(string fileName);

void SortListOfStringInstancesInternal(vector<StringInstance>& strings);
vector<StringInstance> SortListOfStringInstances(vector<StringInstance> strings);

// Get (Augmented) Prefix Tree Acceptor from list of Strings
DFA GetPTAFromListOfStringInstances(vector<StringInstance>& strings, bool APTA);