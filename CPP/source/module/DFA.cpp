#include "DFA.h"
#include <pybind11/pybind11.h>
#include <pybind11/stl.h>
#include <sstream>
#include <fstream>
#include <windows.h>
#include <ppl.h>
#include <mutex>

using concurrency::parallel_for_each;
std::mutex m;

State::State(StateStatus stateStatus, unsigned int stateID)
    : stateStatus(stateStatus), stateID(stateID) {}

DFA::DFA(map<unsigned int, State>& states, State& startingState, vector<char>& alphabet)
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
    //SortListOfStringInstancesInternal(strings);

    bool exists;
    unsigned int count;
    vector<char> alphabet;
    map<unsigned int, State> states;
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

PYBIND11_MODULE(DFA_Toolkit, module)
{
    module.doc() = R"pbdoc(
        A DFA-Toolkit
        -----------------------

        .. currentmodule:: DFA_Toolkit

        .. autosummary::
           :toctree: _generate

           StateStatus
           State
           DFA
           DFA.getAcceptingStates
           DFA.addState
           DFA.depth
           DFA.describe
           StringInstance
           GetListOfStringInstancesFromFile
           SortListOfStringInstances
           GetPTAFromListOfStringInstances
           StringInstanceConsistentWithDFA
           ListOfStringInstancesConsistentWithDFA
           GetStringStatusInRegardToDFA
    )pbdoc";

    pybind11::enum_<StateStatus>(module, "StateStatus", "Represents a DFA's state's status. (Accepting/Rejecting/Unknown).")
        .value("ACCEPTING", StateStatus::ACCEPTING, "State is an accepting state.")
        .value("REJECTING", StateStatus::REJECTING, "State is a rejecting state.")
        .value("UNKNOWN", StateStatus::UNKNOWN, "State is neither an accepting nor a rejecting state.")
        .export_values();
    ;

    pybind11::class_<State>(module, "State", "Represents a DFA's state.")
        .def(pybind11::init<StateStatus, unsigned int>(), "constructor", pybind11::arg("stateStatus"), pybind11::arg("stateID"))
        .def_readwrite("stateStatus", &State::stateStatus, "State's status (Accepting/Rejecting/Unknown).")
        .def_readwrite("stateID", &State::stateID, "State's identification number.");

    pybind11::class_<DFA>(module, "DFA", "Represents a DFA.")
        .def(pybind11::init<map<unsigned int, State>&, State&, vector<char>&>(), "constructor", pybind11::arg("states"), pybind11::arg("startingState"), pybind11::arg("alphabet"))
        .def("getAcceptingStates", &DFA::getAcceptingStates, R"pbdoc(
        Returns DFA's accepting states as a list of State objects.

        This method takes no arguments.
    )pbdoc")
        .def("addState", &DFA::addState, R"pbdoc(
        Adds a State object to the DFA's states.

        Takes a StateStatus enum value and an integer for the StateID as arguments. This method does not return anything.
    )pbdoc")
        .def("depth", &DFA::depth, R"pbdoc(
        Returns the DFA's depth.

        Returns the DFA's depth by traversing the DFA. This method does not take any arguments.
    )pbdoc")
        .def("describe", &DFA::describe, R"pbdoc(
        Prints the DFA's details.

        If the boolean argument is true, all of the DFA's details are printed while if it is false, only an overwiew is printed. This method does not return anything.
    )pbdoc")
        .def_readwrite("states", &DFA::states, "DFA's states as a map of the stateIDs to the State objects.")
        .def_readwrite("startingState", &DFA::startingState, "DFA's starting state as a State object.")
        .def_readwrite("alphabet", &DFA::alphabet, "DFA's alphabet as a list of characters.");

    pybind11::class_<StringInstance>(module, "StringInstance", "Represents either a positive, negative or an unknown string instance of a given DFA.")
        .def(pybind11::init<string&, StateStatus, unsigned int&>(), "constructor1", pybind11::arg("text"), pybind11::arg("stringStatus"), pybind11::arg("length"))
        .def(pybind11::init<string&, const string&>(), "constructor2", pybind11::arg("text"), pybind11::arg("delimiter"))
        .def_readwrite("stringStatus", &StringInstance::stringStatus, "String is either an accepting, rejecting or unknown string instance.")
        .def_readwrite("length", &StringInstance::length, "String's length.")
        .def_readwrite("stringValue", &StringInstance::stringValue, "String's value.");

    module.def("GetListOfStringInstancesFromFile", &GetListOfStringInstancesFromFile, R"pbdoc(
        Parses an Abbadingo DFA dataset into a list of StringInstance objects.

        File format should follow Abaddingo dataset structure. Single parameter for file dir/name.
    )pbdoc");

    module.def("SortListOfStringInstances", &SortListOfStringInstances, R"pbdoc(
        Sorts a list of String Instances by length.

    )pbdoc");

    module.def("GetPTAFromListOfStringInstances", &GetPTAFromListOfStringInstances, R"pbdoc(
        Parses a list of StringInstance objects into a APTA or PTA as a DFA object.

        Gets the Augumented Prefix Tree Acceptor or the Prefix Tree Acceptor.
        First Parameter -> List of String Instances
        Second Parameter -> Boolean value (True for APTA and False for PTA)
    )pbdoc");

    module.def("StringInstanceConsistentWithDFA", &StringInstanceConsistentWithDFA, R"pbdoc(
        Checks if a given string instance is consistent with the given DFA.

        Returns a boolean value. True if string instance is consistent with the DFA or vica versa for false.
        First Parameter -> String Instance
        Second Parameter -> DFA Instance
    )pbdoc");

    module.def("ListOfStringInstancesConsistentWithDFA", &ListOfStringInstancesConsistentWithDFA, R"pbdoc(
        Checks if a given list of string instances is consistent with the given DFA.

        Returns a boolean value. True if all string instances within the list are consistent with the DFA or false if at least a single string instance within the list is inconsistent.
        First Parameter -> List of String Instances
        Second Parameter -> DFA Instance
    )pbdoc");

    module.def("GetStringStatusInRegardToDFA", &GetStringStatusInRegardToDFA, R"pbdoc(
        Gives the string status of a given string instance in regard to a given DFA.

        Returns a StateStatus value. This represents either an accepting, rejecting or unknown for the given string instance.
        Please note that the given string's status is irrelevant to this function.
        First Parameter -> String Instance
        Second Parameter -> DFA Instance
    )pbdoc");
}