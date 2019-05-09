* BIT_RATE - 16/32
* in order to simulate 44100 sample rate, a sample should be written every 22.675 micro sec. since sleeping for 22.675 microsec and writing a sample is not accurate
the algorithm should be fill buffer of X samples and then sleep for X*22.675.
* BUF_SIZE_SAMPLES - buffer size to fill oscillator output, in samples.
* registering an osc - curl ec2-34-246-193-176.eu-west-1.compute.amazonaws.com:9876/register/sin
* changing note - curl ec2-34-246-193-176.eu-west-1.compute.amazonaws.com:9876/token/{token}/note/f