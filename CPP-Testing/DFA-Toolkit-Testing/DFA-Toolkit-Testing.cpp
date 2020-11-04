#include <chrono> 
#include "DFA.h"
using namespace std::chrono;

int main()
{
    auto start = high_resolution_clock::now();
    vector<StringInstance> listOfStrings;
    try 
    {
        listOfStrings = GetListOfStringInstancesFromFile("dataset3\\train.a");
    }
    catch (_exception e) {
        std::cout << "Error!" << std::endl;
    }
    //SortListOfStringInstances(listOfStrings);
    DFA APTA = GetPTAFromListOfStringInstances(listOfStrings, true);
    APTA.describe(false);
    auto stop = high_resolution_clock::now();
    auto duration = duration_cast<milliseconds>(stop - start);
    std::cout << "Average time: " << duration.count() << std::endl;
}