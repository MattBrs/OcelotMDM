#pragma once

#include <atomic>
#include <condition_variable>
#include <functional>
#include <mutex>
#include <thread>

namespace OcelotMDM::component {
class Timer {
   public:
    Timer() = default;
    ~Timer();

    void start(std::function<void()> fn, const int intr, const bool cycle);
    void stop();

   private:
    std::thread             th;
    std::condition_variable cv;
    std::mutex              mtx;
    std::atomic<bool>       running;
};
}  // namespace OcelotMDM::component
