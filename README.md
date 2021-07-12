	.
	├── .github
		└── workflows			# Workflow file used in Github Actions.
	├── core        			# Source files within the dfalearningtoolkit module.
	├── datasets    			# Datasets used throughout this project.
		├── Abbadingo   		# Some datasets from the Abbadingo competition (http://abbadingo.cs.nuim.ie/data-sets.html).
		├── Comparison    		# Datasets used to evaluate the correctness of our implementation.
		├── Stamina        		# Datasets within the Stamina competition (http://stamina.chefbe.net/download).
		├── Visualisation   	# Output generated from the visualisation unit tests. 
		├── GI-learning        	# Datases used to evaluate the performance of GI-learning.
		├── Generated Abbadingo	# Datasets generating using the Abbadingo protocol.
		├── LearnLib      		# Datases used to evaluate the performance of LearnLib.
		└── TestingAPTAs		# Datasets used in some of the unit tests. 
	├── doc        				# Documentation in HTML format.
	├── test        			# Source files for the unit tests and the benchmarks.
	├── util        			# Source files within the dfalearningtoolkit.util module.
	├── .gitignore				# .gitignore file.
	├── go.mod      			# Module file 1.
	├── go.sum      			# Module file 2.
	└── main.go					# Main executable with a simple CLI.

Please note that some datasets were ommited from the submission due to their size.
These can be found at: https://github.com/Cherrett/DFA-Learning-Datasets.