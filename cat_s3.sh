#!/bin/bash

parallel 'aws s3 cp {} -; echo ""'
