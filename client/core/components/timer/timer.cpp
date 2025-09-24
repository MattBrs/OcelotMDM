#include "timer.hpp"

#include <chrono>
#include <functional>
#include <iostream>
#include <mutex>
#include <thread>

namespace OcelotMDM::component {
void Timer::start(std::function<void()> fn, const int intr, const bool cycle) {
    if (running.load()) {
        return;
    }

    this->running.store(true);
    this->th = std::thread([this, fn, intr, cycle]() {
        while (this->running.load()) {
            std::unique_lock<std::mutex> lock(this->mtx);
            this->cv.wait_for(lock, std::chrono::milliseconds(intr), [this]() {
                return !this->running.load();
            });

            fn();

            if (!cycle) {
                break;
            }
        }
    });
}

void Timer::stop() {
    this->running.store(false);
    this->cv.notify_one();

    if (this->th.joinable()) {
        this->th.join();
    }
}

Timer::~Timer() {
    this->stop();
}
}  // namespace OcelotMDM::component
