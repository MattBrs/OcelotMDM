### How to build

From the source directory (where there is the main cmake file)
~~~
conan install . --output-folder=build --build=missing
~~~

To generate the build steps
~~~
cd build
cmake .. -GNinja  -DCMAKE_TOOLCHAIN_FILE=build/Release/generators/conan_toolchain.cmake -DCMAKE_BUILD_TYPE=Release
~~~

To build
~~~
ninja
~~~
