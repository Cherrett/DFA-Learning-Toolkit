#include "DFA.h"
#include <sstream>
#include <fstream>
#include <windows.h>
#include <ppl.h>
#include <mutex>

using concurrency::parallel_for_each;
std::mutex m;

State::State(StateStatus stateStatus, unsigned int stateID)
    : stateStatus(stateStatus), stateID(stateID) {}

State::State(StateStatus stateStatus, unsigned int stateID, map<char, unsigned int> transitions)
    : stateStatus(stateStatus), stateID(stateID), transitions(transitions) {}

NFA_State::NFA_State(StateStatus stateStatus, unsigned int stateID)
    : stateStatus(stateStatus), stateID(stateID) {}

NFA_State::NFA_State(StateStatus stateStatus, unsigned int stateID, map<char, vector<unsigned int>> transitions)
    : stateStatus(stateStatus), stateID(stateID), transitions(transitions) {}

//TransitionFunction::TransitionFunction(State& fromState, State& toState, char& symbol)
//    : fromState(fromState), toState(toState), symbol(symbol) {}

DFA::DFA(map<unsigned int, State>& states, State& startingState, vector<char>& alphabet)
    : states(states), startingState(startingState), alphabet(alphabet) {}

NFA::NFA(map<unsigned int, NFA_State>& states, NFA_State& startingState, vector<char>& alphabet)
    : states(states), startingState(startingState), alphabet(alphabet) {}

vector<State> DFA::getAcceptingStates() {
    vector<State> acceptingStates;

    map<unsigned int, State>::iterator it;
    for (it = this->states.begin(); it != this->states.end(); it++)
        if (it->second.stateStatus == StateStatus::ACCEPTING)
            acceptingStates.push_back(it->second);

    return acceptingStates;
}

vector<State> DFA::getRejectingStates() {
    vector<State> rejectingStates;

    map<unsigned int, State>::iterator it;
    for (it = this->states.begin(); it != this->states.end(); it++)
        if (it->second.stateStatus == StateStatus::REJECTING)
            rejectingStates.push_back(it->second);

    return rejectingStates;
}

void DFA::addState(StateStatus& stateStatus) {
    this->states[this->states.size()] = State(stateStatus, this->states.size());
}

//void DFA::addState(State& state) {
//    this->states[state.stateID] = state;
//}

//void DFA::addTransitionFunction(State& fromState, State& toState, char& symbol) {
//    this->transitionFunctions.emplace_back(fromState, toState, symbol);
//}

unsigned int DFA::depth() {
    map<unsigned int, unsigned int> stateMap;

    depthUtil(this->startingState, 0, stateMap);

    unsigned int max_value = 0;
    std::map<unsigned int, unsigned int>::iterator map_iterator;
    for (map_iterator = stateMap.begin(); map_iterator != stateMap.end(); ++map_iterator) {
        if (map_iterator->second > max_value)
            max_value = map_iterator->second;
    }
    return max_value;
}

void DFA::depthUtil(State& state, int count, map<unsigned int, unsigned int>& stateMap) {
    stateMap[state.stateID] = count;

    std::map<char, unsigned int>::iterator transitions_iterator;
    for (transitions_iterator = state.transitions.begin(); transitions_iterator != state.transitions.end(); ++transitions_iterator) {
        if (stateMap.count(transitions_iterator->second) == 0)
            depthUtil(this->states[transitions_iterator->second], count + 1, stateMap);
    }

    /*for (TransitionFunction& transitionFunction : this->transitionFunctions) {
        if (transitionFunction.fromState.stateID == stateID && stateMap.count(transitionFunction.toState.stateID) == 0) {
            depthUtil(transitionFunction.toState.stateID, count + 1, stateMap);
        }
    }*/
}

void DFA::describe(bool detail) {
    std::cout << "This DFA has " << this->states.size() << " states and " << this->alphabet.size() << " alphabet" << std::endl;
    if (detail) {
        std::cout << "States:" << std::endl;
        //for (State& state : this->states) {

        map<unsigned int, State>::iterator statesIterator;
        for (statesIterator = this->states.begin(); statesIterator != this->states.end(); statesIterator++)
        {
            if (statesIterator->second.stateStatus == StateStatus::ACCEPTING) {
                std::cout << statesIterator->first << " ACCEPTING" << std::endl;
            }
            else if (statesIterator->second.stateStatus == StateStatus::REJECTING) {
                std::cout << statesIterator->first << " REJECTING" << std::endl;
            }
            else {
                std::cout << statesIterator->first << " UNKNOWN" << std::endl;
            }
        }
        std::cout << "Accepting States:" << std::endl;
        for (State& state : this->getAcceptingStates()) {
            std::cout << state.stateID << std::endl;
        }
        std::cout << "Starting State:" << std::endl << this->startingState.stateID << std::endl;
        std::cout << "Alphabet:" << std::endl;
        for (char& character : this->alphabet) {
            std::cout << character << std::endl;
        }
        /*std::cout << "Transition Functions:" << std::endl;
        for (TransitionFunction& transitionFunction : this->transitionFunctions) {
            std::cout << transitionFunction.fromState.stateID << "->" << transitionFunction.toState.stateID << "=" << transitionFunction.symbol << std::endl;
        }*/
    }
}

StringInstance::StringInstance(string& stringValue, StateStatus stringStatus, unsigned int& length)
    : stringValue(stringValue), stringStatus(stringStatus), length(length) {}

StringInstance::StringInstance(string& text, const string& delimiter) {
    this->stringValue = "";
    size_t pos = 0;
    std::string token;

    // accepting / rejecting
    pos = text.find(delimiter);
    token = text.substr(0, pos);

    if (token == "0")
        this->stringStatus = StateStatus::REJECTING;
    else if (token == "1")
        this->stringStatus = StateStatus::ACCEPTING;
    else if (token == "-1")
        this->stringStatus = StateStatus::UNKNOWN;
    else
        throw "Error, unkwown string status. Value: '" + token + "' .";

    text.erase(0, pos + 1);

    // length
    pos = text.find(delimiter);

    if (pos != std::string::npos) {
        token = text.substr(0, pos);
        std::istringstream(token) >> this->length;
        text.erase(0, pos + 1);

        while ((pos = text.find(delimiter)) != std::string::npos) {
            token = text.substr(0, pos);
            this->stringValue.append(token);
            text.erase(0, pos + 1);
        }

        if (text.length() != 0) {
            this->stringValue.append(text);
        }
    }
    else {
        this->length = 0;
    }
}

bool StringInstance::operator< (const StringInstance& otherString) const {
    return length < otherString.length;
}

vector<StringInstance> GetListOfStringInstancesFromFile(string fileName) {
    vector<StringInstance> listOfStrings;
    std::ifstream infile(fileName);
    string line;
    // ignore first line
    std::getline(infile, line);
    if (line.length() == 0) {
        throw "Error, Invalid file name";
    }
    while (std::getline(infile, line))
    {
        listOfStrings.emplace_back(line, " ");
    }
    return listOfStrings;
}

void SortListOfStringInstancesInternal(vector<StringInstance>& strings) {
    sort(strings.begin(), strings.end());
}

vector<StringInstance> SortListOfStringInstances(vector<StringInstance> strings) {
    sort(strings.begin(), strings.end());
    return strings;
}

DFA GetPTAFromListOfStringInstances(vector<StringInstance>& strings, bool APTA) {
    SortListOfStringInstancesInternal(strings);

    bool exists;
    unsigned int count;
    vector<char> alphabet;
    map<unsigned int, State> states;
    //vector<TransitionFunction> transitionFunctions;
    unsigned int startingStateID = 0;
    unsigned int currentStateID;

    if (strings[0].length == 0) {
        if (strings[0].stringStatus == StateStatus::ACCEPTING) {
            states[0] = State(StateStatus::ACCEPTING, 0);
        }
        else {
            states[0] = State(StateStatus::REJECTING, 0);
        }
    }
    else {
        states[0] = State(StateStatus::UNKNOWN, 0);
    }

    for (StringInstance& string : strings) {
        if (!APTA && string.stringStatus != StateStatus::ACCEPTING)
            continue;
        currentStateID = startingStateID;
        count = 0;
        for (char& character : string.stringValue) {
            count++;
            exists = false;
            // alphabet check
            if (std::find(alphabet.begin(), alphabet.end(), character) == alphabet.end())
                alphabet.push_back(character);

            map<char, unsigned int>::iterator transitionIterator = states[currentStateID].transitions.find(character);
            if (transitionIterator != states[currentStateID].transitions.end())
            {
                currentStateID = transitionIterator->second;
                exists = true;
            }

            /*for (TransitionFunction& transitionFunction : transitionFunctions) {
                if (transitionFunction.fromState.stateID == currentState.stateID && transitionFunction.symbol == character) {
                    currentState = transitionFunction.toState;
                    exists = true;
                    break;
                }
            }*/

            if (!exists) {
                // last symbol in string check
                if (count == string.stringValue.size()) {
                    if (string.stringStatus == StateStatus::ACCEPTING)
                        states[states.size()] = State(StateStatus::ACCEPTING, states.size());
                    else
                        states[states.size()] = State(StateStatus::REJECTING, states.size());
                }
                else {
                    states[states.size()] = State(StateStatus::UNKNOWN, states.size());
                }
                states[currentStateID].transitions[character] = states[states.size() - 1].stateID;
                //transitionFunctions.emplace_back(currentState, states[states.size() - 1], character);
                currentStateID = states[states.size() - 1].stateID;
            }
            else {
                // last symbol in string check
                if (count == string.stringValue.size()) {
                    if (string.stringStatus == StateStatus::ACCEPTING) {
                        if (states[currentStateID].stateStatus == StateStatus::REJECTING)
                            throw "Error, state already set to rejecting, cannot set to accepting";
                        else
                            states[currentStateID].stateStatus = StateStatus::ACCEPTING;
                    }
                    else {
                        if (states[currentStateID].stateStatus == StateStatus::ACCEPTING)
                            throw "Error, state already set to accepting, cannot set to rejecting";
                        else
                            states[currentStateID].stateStatus = StateStatus::REJECTING;
                    }
                }
            }
        }
    }

    return DFA(states, states[startingStateID], alphabet);
}

bool StringInstanceConsistentWithDFA(StringInstance& string, DFA& dfa) {
    // Skip unknown strings (test data)
    if (string.stringStatus == StateStatus::UNKNOWN) {
        return true;
    }

    bool exists;
    State currentState = dfa.startingState;
    unsigned int count = 0;
    for (char& character : string.stringValue) {
        count++;
        exists = false;

        map<char, unsigned int>::iterator transitionIterator = currentState.transitions.find(character);
        if (transitionIterator != currentState.transitions.end())
        {
            currentState = dfa.states[transitionIterator->second];
            exists = true;
        }

        if (!exists) {
            return false;
        }
        else {
            // last symbol in string check
            if (count == string.length) {
                if (string.stringStatus == StateStatus::ACCEPTING) {
                    if (currentState.stateStatus == StateStatus::REJECTING) {
                        return false;
                    }
                }
                else {
                    if (currentState.stateStatus == StateStatus::ACCEPTING) {
                        return false;
                    }
                }
            }
        }
    }
    return true;
}

bool ListOfStringInstancesConsistentWithDFA(vector<StringInstance>& strings, DFA& dfa) {
    bool consistent = true;

    parallel_for_each(begin(strings), end(strings), [&](StringInstance string) {

        if (!consistent) {
            return;
        }
        else {
            if (!StringInstanceConsistentWithDFA(string, dfa)) {
                m.lock();
                consistent = false;
                m.unlock();
                return;
            }
        }
        });

    return consistent;
}

StateStatus GetStringStatusInRegardToDFA(StringInstance& string, DFA& dfa) {
    bool exists;
    State currentState = dfa.startingState;
    unsigned int count = 0;
    for (char& character : string.stringValue) {
        count++;
        exists = false;

        map<char, unsigned int>::iterator transitionIterator = currentState.transitions.find(character);
        if (transitionIterator != currentState.transitions.end())
        {
            currentState = dfa.states[transitionIterator->second];
            exists = true;
        }

        if (!exists) {
            return StateStatus::UNKNOWN;
        }
        else {
            // last symbol in string check
            if (count == string.stringValue.size()) {
                switch (currentState.stateStatus) {
                case StateStatus::ACCEPTING:
                    return StateStatus::ACCEPTING;
                case StateStatus::REJECTING:
                    return StateStatus::REJECTING;
                case StateStatus::UNKNOWN:
                    return StateStatus::UNKNOWN;
                }
            }
        }
    }
    return StateStatus::UNKNOWN;
}

vector<StringInstance> GetAcceptingStringInstances(vector<StringInstance> strings) {
    vector<StringInstance> acceptingStrings;

    for (StringInstance& stringInstance : strings) {
        if (stringInstance.stringStatus == StateStatus::ACCEPTING) {
            acceptingStrings.emplace_back(stringInstance);
        }
    }
    return acceptingStrings;
}

vector<StringInstance> GetRejectingStringInstances(vector<StringInstance> strings) {
    vector<StringInstance> rejectingStrings;

    for (StringInstance& stringInstance : strings) {
        if (stringInstance.stringStatus == StateStatus::REJECTING) {
            rejectingStrings.emplace_back(stringInstance);
        }
    }
    return rejectingStrings;
}

bool isNFAdeterministic(NFA nfa) {
    for (map<unsigned int, NFA_State>::iterator it = nfa.states.begin(); it != nfa.states.end(); it++) {
        for (map<char, vector<unsigned int>>::iterator it2 = it->second.transitions.begin(); it2 != it->second.transitions.end(); it2++) {
            if (it2->second.size() > 1) {
                return false;
            }
        }
    }
    return true;
}

NFA RPNI_Derive(DFA dfa, vector<vector<unsigned int>> partition) {
    map<unsigned int, unsigned int> new_mappings;
    map<unsigned int, NFA_State> states;
    
    for (vector<unsigned int>& currentBlock : partition) {
        StateStatus stateStatus = StateStatus::UNKNOWN;
        map<char, vector<unsigned int>> transitions;

        for (unsigned int& stateID : currentBlock) {
            if (dfa.states[stateID].stateStatus == StateStatus::ACCEPTING) {
                stateStatus = StateStatus::ACCEPTING;
            }
            for (map<char, unsigned int>::iterator it = dfa.states[stateID].transitions.begin(); it != dfa.states[stateID].transitions.end(); it++) {
                if (transitions.count(it->first) == 0) {
                    transitions[it->first] = vector<unsigned int>{it->second};
                }
                else {
                    if (std::find(transitions[it->first].begin(), transitions[it->first].end(), it->second) == transitions[it->first].end()) 
                        transitions[it->first].emplace_back(it->second);
                }
            }
            new_mappings[stateID] = states.size();
        }
        states[states.size()] = NFA_State(stateStatus, states.size(), transitions);
    }
    // update new states via mappings
    for (map<unsigned int, NFA_State>::iterator it = states.begin(); it != states.end(); it++) {
        for (map<char, vector<unsigned int>>::iterator it2 = it->second.transitions.begin(); it2 != it->second.transitions.end(); it2++) {
            for (unsigned int i = 0; i < it2->second.size(); i++)
                it2->second[i] = new_mappings[it2->second[i]];
        }
    }

    return NFA(states, states[new_mappings[dfa.startingState.stateID]], dfa.alphabet);
}

NFA RPNI_Merge(NFA nfa, unsigned int state1, unsigned int state2) {
    StateStatus stateStatus = StateStatus::UNKNOWN;
    for (map<unsigned int, NFA_State>::iterator it = nfa.states.begin(); it != nfa.states.end(); it++) {
        if (it->first == state1) {
            if (it->second.stateStatus == StateStatus::ACCEPTING)
                stateStatus = StateStatus::ACCEPTING;
        }
        else if (it->first == state2) {
            if (it->second.stateStatus == StateStatus::ACCEPTING)
                stateStatus = StateStatus::ACCEPTING;
            for (map<char, vector<unsigned int>>::iterator it2 = it->second.transitions.begin(); it2 != it->second.transitions.end(); it2++) {
                for (unsigned int i = 0; i < it2->second.size(); i++) {
                    if (nfa.states[state1].transitions.count(it2->first) == 0) {
                        nfa.states[state1].transitions[it2->first] = vector<unsigned int>{ it2->second.at(i) };
                    }
                    else {
                        if (std::find(nfa.states[state1].transitions[it2->first].begin(), nfa.states[state1].transitions[it2->first].end(), it2->second.at(i)) == nfa.states[state1].transitions[it2->first].end())
                            nfa.states[state1].transitions[it2->first].emplace_back(it2->second.at(i));               
                    }
                }
            }
            continue;
        }
        
        for (map<char, vector<unsigned int>>::iterator it2 = it->second.transitions.begin(); it2 != it->second.transitions.end(); it2++) {
            for (unsigned int i = 0; i < it2->second.size(); i++)
                if (it2->second[i] == state2)
                    if (std::find(it2->second.begin(), it2->second.end(), state1) == it2->second.end())
                        it2->second[i] = state1;
                    else
                        it2->second.erase(it2->second.begin() + i);
        }
    }

    nfa.states.erase(state2);
    nfa.states[state1].stateStatus = stateStatus;

    if (nfa.startingState.stateID == state2)
        nfa.startingState = nfa.states[state1];

    return nfa;
}

RPNI_Deterministic_Merge_object::RPNI_Deterministic_Merge_object(vector<vector<unsigned int>> partition, DFA dfa)
    : partition(partition), dfa(dfa) {};

RPNI_Deterministic_Merge_object RPNI_Deterministic_Merge(NFA nfa, vector<vector<unsigned int>> partition) {
    //map<unsigned int, State> newStates;
    
    while (!isNFAdeterministic(nfa)) {
        //vector<vector<unsigned int>> newPartition;
        bool exit_loop = false;
        for (map<unsigned int, NFA_State>::iterator it = nfa.states.begin(); it != nfa.states.end(); it++) {
            //vector<unsigned int> newBlock;
            for (map<char, vector<unsigned int>>::iterator it2 = it->second.transitions.begin(); it2 != it->second.transitions.end(); it2++) {
                if (it2->second.size() > 1) {
                    // merging
                    unsigned int state1 = it2->second[0];
                    unsigned int state2 = it2->second[1];
                    nfa = RPNI_Merge(nfa, state1, state2);
                    for (unsigned int element : partition.at(state2))
                        partition.at(state1).push_back(element);
                    partition.erase(partition.begin()+state2);
                    // need to create normalization function after merge
                    exit_loop = true;
                    break;
                }
            }
            if (exit_loop)
                break;
        }
    }
    

    return RPNI_Deterministic_Merge_object(partition, NFAtoDFA(nfa));
}

DFA NFAtoDFA(NFA nfa) {
    map<unsigned int, State> states;

    for (map<unsigned int, NFA_State>::iterator it = nfa.states.begin(); it != nfa.states.end(); it++) {
        states[it->first] = State(it->second.stateStatus, it->second.stateID);
        for (map<char, vector<unsigned int>>::iterator it2 = it->second.transitions.begin(); it2 != it->second.transitions.end(); it2++) {
            states[it->first].transitions[it2->first] = it2->second.at(0);
        }
    }    

    return DFA(states, states[nfa.startingState.stateID], nfa.alphabet);
}

bool RPNI_StringInstanceConsistentWithDFA(StringInstance& string, DFA& dfa) {
    bool exists;
    State currentState = dfa.startingState;
    unsigned int count = 0;

    if (string.length == 0) {
        if (dfa.startingState.stateStatus == StateStatus::ACCEPTING) {
            return false;
        }
    }
    else {
        for (char& character : string.stringValue) {
            count++;
            exists = false;

            map<char, unsigned int>::iterator transitionIterator = currentState.transitions.find(character);
            if (transitionIterator != currentState.transitions.end())
            {
                currentState = dfa.states[transitionIterator->second];
                exists = true;
            }

            if (!exists) {
                return true;
            }
            else {
                // last symbol in string check
                if (count == string.stringValue.size()) {
                    if (currentState.stateStatus == StateStatus::ACCEPTING) {
                        return false;
                    }
                }
            }
        }
    }
    return true;
}

bool RPNI_ListOfNegativeStringInstancesConsistentWithDFA(vector<StringInstance>& strings, DFA& dfa) {
    bool consistent = true;

    parallel_for_each(begin(strings), end(strings), [&](StringInstance string) {

        if (!consistent) {
            return;
        }
        else {
            if (!RPNI_StringInstanceConsistentWithDFA(string, dfa)) {
                m.lock();
                consistent = false;
                m.unlock();
                return;
            }
        }
        });

    return consistent;
}

DFA RPNI(vector<StringInstance>& acceptingStrings, vector<StringInstance>& rejectingStrings) {
    DFA PTA = GetPTAFromListOfStringInstances(acceptingStrings, false);
    DFA currentHypothesis = PTA;
    DFA tempHypothesis = PTA;
    vector<vector<unsigned int>> currentPartition;

    for (map<unsigned int, State>::iterator it = currentHypothesis.states.begin(); it != currentHypothesis.states.end(); it++)
        currentPartition.emplace_back(vector<unsigned int>{it->first});

    for (int i = 1; i < PTA.states.size(); i++) {
        for (int j = 0; j < i; j++) {
            // merge the block which contains state i with the block which contains state j
            vector<vector<unsigned int>> tempPartition;
            
            vector<unsigned int> tempBlock = vector<unsigned int>{};
            for (vector<unsigned int> block : currentPartition) {
                if (std::find(block.begin(), block.end(), i) != block.end()) {
                    for (unsigned int state : block)
                        tempBlock.emplace_back(state);
                }
                else if (std::find(block.begin(), block.end(), j) != block.end()){
                    for (unsigned int state : block)
                        tempBlock.emplace_back(state);
                }
                else {
                    tempPartition.emplace_back(block);
                }
            }
            tempPartition.emplace_back(tempBlock);

            // get quotient automaton
            NFA tempHypothesisNFA = RPNI_Derive(PTA, tempPartition);

            if (isNFAdeterministic(tempHypothesisNFA)) {
                tempHypothesis = NFAtoDFA(tempHypothesisNFA);            
            }
            else {
                // determinize the quotient automaton (if necessary) by state merging
                RPNI_Deterministic_Merge_object temp = RPNI_Deterministic_Merge(tempHypothesisNFA, tempPartition);
                tempHypothesis = temp.dfa;
                tempPartition = temp.partition;
            }
            // Does tempDFA reject all negative strings?
            if (RPNI_ListOfNegativeStringInstancesConsistentWithDFA(rejectingStrings, tempHypothesis)) {
                // Treat tempHypothesis as the current hypothesis
                currentHypothesis = tempHypothesis;
                currentPartition = tempPartition;
            }    
        }
    }
    return currentHypothesis;
}