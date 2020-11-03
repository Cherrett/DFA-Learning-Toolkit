# import PreProcessing
# import String
# import DFA
from DFA_Toolkit import GetListOfStringsFromFile, GetPTAFromListOfStrings, UNKNOWN, State
from timeit import Timer


def test():
    listOfStrings = GetListOfStringsFromFile('dataset3\\train.a')
    APTA = GetPTAFromListOfStrings(listOfStrings, True)
    APTA.describe(False)

    temp = State(UNKNOWN, len(APTA.states))
    APTA.states.append(temp)
# APTA = PreProcessing.GetAPTAFromMCA(MCA)
# MCA.describe(False)
# for string in listOfStrings:
#     print(string.string)

# alphabet = [ord('0'), ord('1')]
# states = []
# transitionFunctions = []
# state1 = DFA.createState(1, len(states))
# states.append(state1)
# for i in range(10000):
# state1 = DFA.createState(DFA.ACCEPTING, len(states))
# states.append(state1)
# state2 = DFA.createState(DFA.REJECTING, len(states))
# states.append(state2)
# state3 = DFA.createState(DFA.ACCEPTING, len(states))
# states.append(state3)
# state4 = DFA.createState(DFA.ACCEPTING, len(states))
# states.append(state4)
#
# transitionFunction = DFA.createTransitionFunction(state1, state2, alphabet[1])
# transitionFunctions.append(transitionFunction)
#
# transitionFunction = DFA.createTransitionFunction(state1, state4, alphabet[0])
# transitionFunctions.append(transitionFunction)
#
# transitionFunction = DFA.createTransitionFunction(state2, state3, alphabet[1])
# transitionFunctions.append(transitionFunction)
#
# transitionFunction = DFA.createTransitionFunction(state3, state3, alphabet[1])
# transitionFunctions.append(transitionFunction)
#
# transitionFunction = DFA.createTransitionFunction(state4, state3, alphabet[1])
# transitionFunctions.append(transitionFunction)
#
# transitionFunction = DFA.createTransitionFunction(state4, state4, alphabet[0])
# transitionFunctions.append(transitionFunction)
#
# dfa = DFA.DFA(states, state1, alphabet, transitionFunctions)
# dfa.describe()

# print('Starting state: ', len(dfa.getStartingState()))
# print('number of states: ', len(dfa.getStates()))
# print('number of transitionFunctions: ', len(dfa.getTransitionFunctions()))
# print('started timings')
Cython = Timer(lambda: test())
Cython_value = Cython.timeit(number=1)
print('Average time:', Cython_value)
# Original = Timer(lambda: dfa.describe())
# Original_value = Original.timeit(number=1000)
# print('Original code average time:', Original_value)
# print('Cython is', int(Original_value/Cython_value), 'times faster')
