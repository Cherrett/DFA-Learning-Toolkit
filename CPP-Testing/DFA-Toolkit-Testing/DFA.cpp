#include "DFA.h"
#include <sstream>
#include <fstream>

enum class StateStatus {
    ACCEPTING = 1,
    REJECTING = 0,
    UNKNOWN = 2
};

State::State(StateStatus stateStatus, int stateID)
    : stateStatus(stateStatus), stateID(stateID) {}

TransitionFunction::TransitionFunction(State& fromState, State& toState, char& symbol)
    : fromState(fromState), toState(toState), symbol(symbol) {}

DFA::DFA(vector<State>& states, State& startingState, vector<char>& alphabet, vector<TransitionFunction>& transitionFunctions)
    : states(states), startingState(startingState), alphabet(alphabet), transitionFunctions(transitionFunctions) {}

vector<State> DFA::getAcceptingStates() {
    vector<State> acceptingStates;

    for (State state : this->states)
        if (state.stateStatus == StateStatus::ACCEPTING)
            acceptingStates.push_back(state);
    return acceptingStates;
}

void DFA::addState(StateStatus& stateStatus, unsigned int& statusID) {
    this->states.push_back(State(stateStatus, statusID));
}

void DFA::describe(bool detail) {
    std::cout << "This DFA has " << this->states.size() << " states and " << this->alphabet.size() << " alphabet" << std::endl;
    if (detail) {
        std::cout << "States:" << std::endl;
        for (State state : this->states) {
            if (state.stateStatus == StateStatus::ACCEPTING) {
                std::cout << state.stateID << " ACCEPTING" << std::endl;
            }
            else if (state.stateStatus == StateStatus::REJECTING) {
                std::cout << state.stateID << " REJECTING" << std::endl;
            }
            else {
                std::cout << state.stateID << " UNKNOWN" << std::endl;
            }
        }
        std::cout << "Accepting States:" << std::endl;
        for (State state : this->getAcceptingStates()) {
            std::cout << state.stateID << std::endl;
        }
        std::cout << "Starting State:" << std::endl << this->startingState.stateID << std::endl;
        std::cout << "Alphabet:" << std::endl;
        for (char character : this->alphabet) {
            std::cout << character << std::endl;
        }
        std::cout << "Transition Functions:" << std::endl;
        for (TransitionFunction& transitionFunction : this->transitionFunctions) {
            std::cout << transitionFunction.fromState.stateID << "->" << transitionFunction.toState.stateID << "=" << transitionFunction.symbol << std::endl;
        }
    }
}

String::String(string& stringValue, bool accepting, unsigned int& length)
    : stringValue(stringValue), accepting(accepting), length(length) {}

String::String(string& text, const string& delimiter) {
    this->stringValue = "";
    size_t pos = 0;
    std::string token;

    // accepting / rejecting
    pos = text.find(delimiter);
    token = text.substr(0, pos);
    if (token == "1")
        this->accepting = true;
    else
        this->accepting = false;
    text.erase(0, pos + 1);

    // length
    pos = text.find(delimiter);
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

vector<String> GetListOfStringsFromFile(string fileName) {
    vector<String> listOfStrings;
    std::ifstream infile(fileName);
    string line;
    // ignore first line
    std::getline(infile, line);
    if (line.length() == 0) {
        throw "File not valid";
    }
    while (std::getline(infile, line))
    {
        listOfStrings.push_back(String(line, " "));
    }
    return listOfStrings;
}

DFA GetPTAFromListOfStrings(vector<String>& strings, bool APTA) {
    bool exists;
    unsigned int count;
    vector<char> alphabet;
    vector<State> states;
    vector<TransitionFunction> transitionFunctions;
    State startingState, currentState;

    startingState = State(StateStatus::UNKNOWN, 0);
    states.push_back(startingState);

    for (String& string : strings) {
        if (!APTA && !string.accepting)
            continue;
        currentState = startingState;

        count = 0;
        for (char& character : string.stringValue) {
            count++;
            exists = false;
            // alphabet check
            if (std::find(alphabet.begin(), alphabet.end(), character) == alphabet.end())
                alphabet.push_back(character);

            for (TransitionFunction& transitionFunction : transitionFunctions) {
                if (transitionFunction.fromState.stateID == currentState.stateID && transitionFunction.symbol == character) {
                    currentState = transitionFunction.toState;
                    exists = true;
                    break;
                }
            }

            if (!exists) {
                // last symbol in string check
                if (count == string.stringValue.size()) {
                    if (string.accepting)
                        states.push_back(State(StateStatus::ACCEPTING, static_cast<int>(states.size())));
                    else
                        states.push_back(State(StateStatus::REJECTING, static_cast<int>(states.size())));
                }
                else {
                    states.push_back(State(StateStatus::UNKNOWN, static_cast<int>(states.size())));
                }
                transitionFunctions.push_back(TransitionFunction(currentState, states[states.size() - 1], character));
                currentState = states[states.size() - 1];
            }
            else {
                // last symbol in string check
                if (count == string.stringValue.size()) {
                    if (string.accepting) {
                        if (currentState.stateStatus == StateStatus::REJECTING)
                            throw "Error, state already set to rejecting, cannot set to accepting";
                        else
                            currentState.stateStatus = StateStatus::ACCEPTING;
                    }
                    else {
                        if (currentState.stateStatus == StateStatus::ACCEPTING)
                            throw "Error, state already set to accepting, cannot set to rejecting";
                        else
                            currentState.stateStatus = StateStatus::REJECTING;
                    }
                }
            }
        }
    }

    return DFA(states, startingState, alphabet, transitionFunctions);
}