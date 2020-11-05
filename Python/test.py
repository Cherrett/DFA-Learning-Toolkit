from DFA_Toolkit import GetListOfStringInstancesFromFile, GetPTAFromListOfStringInstances, UNKNOWN, State, SortListOfStringInstances
from timeit import Timer


def test():
    listOfStrings = GetListOfStringInstancesFromFile('dataset3\\train.a')
    # listOfStrings = SortListOfStringInstances(listOfStrings)
    APTA = GetPTAFromListOfStringInstances(listOfStrings, True)
    APTA.describe(False)

    print("DFA Depth:", APTA.depth())

    temp = State(UNKNOWN, len(APTA.states))
    APTA.states.append(temp)

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
