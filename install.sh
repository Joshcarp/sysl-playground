#!/bin/bash

while true; do \
        make main; \
        fswatch -qre close_write .; \
    done 