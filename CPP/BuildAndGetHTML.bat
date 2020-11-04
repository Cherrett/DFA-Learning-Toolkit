pip install ./
cd docs
rmdir /Q /S _build
rmdir /Q /S _generate
make html
cd ..