#include <chrono> 
#include "DFA.h"
using namespace std::chrono;

int main()
{
    auto start = high_resolution_clock::now();
    vector<StringInstance> listOfStrings;
    try 
    {
        listOfStrings = GetListOfStringInstancesFromFile("dataset4\\train.a");
    }
    catch (const char* msg) {
        std::cerr << msg << std::endl;
        exit(-1);
    }
    //SortListOfStringInstances(listOfStrings);
    try
    {
        DFA APTA = GetPTAFromListOfStringInstances(listOfStrings, true);
        APTA.describe(false);
        
        std::cout << "DFA Depth: " << APTA.depth() << std::endl;

        //vector <StringInstance> listOfStringsTesting = GetListOfStringInstancesFromFile("dataset3\\test.a");
        if (ListOfStringInstancesConsistentWithDFA(listOfStrings, APTA)) {
            std::cout << "Consistent" << std::endl;
        }else{
            std::cout << "Not Consistent" << std::endl;
        }

        /*StateStatus stateStatus = GetStringStatusInRegardToDFA(listOfStrings[0], APTA);
        stateStatus = GetStringStatusInRegardToDFA(listOfStrings[8], APTA);
        stateStatus = GetStringStatusInRegardToDFA(listOfStrings[9], APTA);
        listOfStrings = GetListOfStringInstancesFromFile("dataset3\\test.a");
        stateStatus = GetStringStatusInRegardToDFA(listOfStrings[0], APTA);
        std::cout << "Done" << std::endl;*/

        //vector<StringInstance> acceptingStrings = GetAcceptingStringInstances(listOfStrings);
        //vector<StringInstance> rejectingStrings = GetRejectingStringInstances(listOfStrings);

        //DFA finalDFA = RPNI(acceptingStrings, rejectingStrings);
    }
    catch (const char* msg) {
        std::cerr << msg << std::endl;
        exit(-1);
    }
    
    auto stop = high_resolution_clock::now();
    auto duration = duration_cast<milliseconds>(stop - start);
    std::cout << "Time: " << duration.count() << std::endl;
}