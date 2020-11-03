#include "DFA.h"
#include <pybind11/pybind11.h>
#include <pybind11/stl.h>
#include <sstream>
#include <fstream>

enum class StateStatus {
    ACCEPTING = 1,
    REJECTING = 0,
    UNKNOWN = 2
};

State::State(StateStatus stateStatus, unsigned int stateID)
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
    this->states.emplace_back(stateStatus, statusID);
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

StringInstance::StringInstance(string& stringValue, bool accepting, unsigned int& length)
    : stringValue(stringValue), accepting(accepting), length(length) {}

StringInstance::StringInstance(string& text, const string& delimiter) {
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

vector<StringInstance> GetListOfStringsFromFile(string fileName) {
    vector<StringInstance> listOfStrings;
    std::ifstream infile(fileName);
    string line;
    // ignore first line
    std::getline(infile, line);
    while (std::getline(infile, line))
    {
        listOfStrings.emplace_back(line, " ");
    }
    return listOfStrings;
}

DFA GetPTAFromListOfStringInstances(vector<StringInstance>& strings, bool APTA) {
    bool exists;
    unsigned int count;
    vector<char> alphabet;
    vector<State> states;
    vector<TransitionFunction> transitionFunctions;
    State startingState, currentState;

    startingState = State(StateStatus::UNKNOWN, 0);
    states.push_back(startingState);

    for (StringInstance& string : strings) {
        if (!APTA && !string.accepting)
            continue;
        currentState = startingState;

        if (string.length == 0) {
            if (string.accepting) {
                if (startingState.stateStatus == StateStatus::REJECTING) {
                    throw "Error, starting state already set to rejecting, cannot set to accepting.";
                } 
                else {
                    startingState.stateStatus = StateStatus::ACCEPTING;
                }
            }
            else{
                if (startingState.stateStatus == StateStatus::ACCEPTING) {
                    throw "Error, starting state already set to accepting, cannot set to rejecting.";
                }
                else {
                    startingState.stateStatus = StateStatus::REJECTING;
                }
            }
        }
        else {
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
                            states.emplace_back(StateStatus::ACCEPTING, static_cast<unsigned int>(states.size()));
                        else
                            states.emplace_back(StateStatus::REJECTING, static_cast<unsigned int>(states.size()));
                    }
                    else {
                        states.emplace_back(StateStatus::UNKNOWN, static_cast<unsigned int>(states.size()));
                    }
                    transitionFunctions.emplace_back(currentState, states[states.size() - 1], character);
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
    }

    return DFA(states, startingState, alphabet, transitionFunctions);
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
           TransitionFunction
           DFA
		   DFA.getAcceptingStates
		   DFA.addState
		   DFA.describe
           StringInstance
           GetListOfStringsFromFile
           GetPTAFromListOfStringInstances 
    )pbdoc";
    module.def("GetListOfStringsFromFile", &GetListOfStringsFromFile, R"pbdoc(
        Parses an Abbadingo DFA dataset into a list of StringInstance objects.

        File format should follow Abaddingo dataset structure. Single parameter for file dir/name.
    )pbdoc");
    module.def("GetPTAFromListOfStringInstances", &GetPTAFromListOfStringInstances, R"pbdoc(
        Parses a list of StringInstance objects into a APTA or PTA as a DFA object.

        Gets the Augumented Prefix Tree Acceptor or the Prefix Tree Acceptor.
        First Parameter -> List of Strings
        Second Parameter -> Boolean value (True for APTA and False for PTA)
    )pbdoc");

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

    pybind11::class_<TransitionFunction>(module, "TransitionFunction", "Represents a DFA's transition function.")
        .def(pybind11::init<State&, State&, char&>(), "constructor", pybind11::arg("fromState"), pybind11::arg("toState"), pybind11::arg("symbol"))
        .def_readwrite("fromState", &TransitionFunction::fromState, "Transition function's start state.")
        .def_readwrite("toState", &TransitionFunction::toState, "Transition function's end state.")
        .def_readwrite("symbol", &TransitionFunction::symbol, "Transition function's symbol (from alphabet).");

    pybind11::class_<DFA>(module, "DFA", "Represents a DFA.")
        .def(pybind11::init<vector<State>&, State&, vector<char>&, vector<TransitionFunction>&>(), "constructor", pybind11::arg("states"), pybind11::arg("startingState"), pybind11::arg("alphabet"), pybind11::arg("transitionFunctions"))
        .def("getAcceptingStates", &DFA::getAcceptingStates, R"pbdoc(
        Returns DFA's accepting states as a list of State objects.

        This method takes no arguments.
    )pbdoc")
        .def("addState", &DFA::addState, R"pbdoc(
        Adds a State object to the DFA's states.

        Takes a StateStatus enum value and an integer for the StateID as arguments. This method does not return anything.
    )pbdoc")
        .def("describe", &DFA::describe, R"pbdoc(
        Prints DFA's details.

        If the boolean argument is true, all of the DFA's details are printed while if it is false, only an overwiew is printed. This method does not return anything.
    )pbdoc")
        .def_readwrite("states", &DFA::states, "DFA's states as a list of State objects.")
        .def_readwrite("startingState", &DFA::startingState, "DFA's starting state as a State object.")
        .def_readwrite("alphabet", &DFA::alphabet, "DFA's alphabet as a list of characters.")
        .def_readwrite("transitionFunctions", &DFA::transitionFunctions, "DFA's transition functions as a list of TransitionFunction objects.");

    pybind11::class_<StringInstance>(module, "StringInstance", "Represents either a positive or a negative string instance of a given DFA.")
        .def(pybind11::init<string&, bool, unsigned int&>(), "constructor1", pybind11::arg("text"), pybind11::arg("accepting"), pybind11::arg("length"))
        .def(pybind11::init<string&, const string&>(), "constructor2", pybind11::arg("text"), pybind11::arg("delimiter"))
        .def_readwrite("accepting", &StringInstance::accepting, "String is an accepting string if value is true and vica versa for false.")
        .def_readwrite("length", &StringInstance::length, "String's length.")
        .def_readwrite("stringValue", &StringInstance::stringValue, "String's value.");
}