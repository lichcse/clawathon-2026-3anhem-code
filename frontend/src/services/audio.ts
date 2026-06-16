export class AudioService {
  private audioContext: AudioContext | null = null;
  private mediaStreamAudioSourceNode: MediaStreamAudioSourceNode | null = null;
  private scriptProcessor: ScriptProcessorNode | null = null;
  private localStream: MediaStream | null = null;
  private remoteAudioElements: Map<string, HTMLAudioElement> = new Map();
  private remoteStreams: Map<string, MediaStream> = new Map();
  private playbackContext: AudioContext | null = null;
  private playbackScheduleTime: Map<string, number> = new Map();

  onAudioData: ((base64: string) => void) | null = null;

  async startCapture(): Promise<MediaStream> {
    try {
      this.localStream = await navigator.mediaDevices.getUserMedia({
        audio: {
          echoCancellation: true,
          noiseSuppression: true,
          autoGainControl: true,
        },
      });

      this.audioContext = new (window.AudioContext || (window as any).webkitAudioContext)();
      this.mediaStreamAudioSourceNode = this.audioContext.createMediaStreamSource(this.localStream);
      this.scriptProcessor = this.audioContext.createScriptProcessor(4096, 1, 1);

      this.scriptProcessor.onaudioprocess = (event) => {
        if (!this.onAudioData) return;
        const track = this.localStream?.getAudioTracks()[0];
        if (!track?.enabled) return;

        const inputData = event.inputBuffer.getChannelData(0);
        const rms = Math.sqrt(inputData.reduce((sum, s) => sum + s * s, 0) / inputData.length);
        if (rms < 0.01) return; // skip silence frames

        this.onAudioData(this.float32ToBase64(inputData));
      };

      this.mediaStreamAudioSourceNode.connect(this.scriptProcessor);
      this.scriptProcessor.connect(this.audioContext.destination);

      return this.localStream;
    } catch (error) {
      console.error('Error accessing microphone:', error);
      throw error;
    }
  }

  async stopCapture() {
    this.onAudioData = null;

    if (this.localStream) {
      this.localStream.getTracks().forEach((track) => track.stop());
      this.localStream = null;
    }

    if (this.scriptProcessor) {
      this.scriptProcessor.disconnect();
      this.scriptProcessor = null;
    }

    if (this.mediaStreamAudioSourceNode) {
      this.mediaStreamAudioSourceNode.disconnect();
      this.mediaStreamAudioSourceNode = null;
    }

    if (this.audioContext) {
      await this.audioContext.close();
      this.audioContext = null;
    }
  }

  setMicMuted(isMuted: boolean) {
    if (this.localStream) {
      this.localStream.getAudioTracks().forEach((track) => {
        track.enabled = !isMuted;
      });
    }
  }

  playRemoteAudio(userId: string, base64: string) {
    if (!this.playbackContext) {
      this.playbackContext = new AudioContext();
    }
    if (this.playbackContext.state === 'suspended') {
      this.playbackContext.resume();
    }

    const binary = atob(base64);
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
      bytes[i] = binary.charCodeAt(i);
    }
    const float32 = new Float32Array(bytes.buffer);

    const ctx = this.playbackContext;
    const audioBuffer = ctx.createBuffer(1, float32.length, ctx.sampleRate);
    audioBuffer.copyToChannel(float32, 0);

    const source = ctx.createBufferSource();
    source.buffer = audioBuffer;
    source.connect(ctx.destination);

    const now = ctx.currentTime;
    const scheduled = this.playbackScheduleTime.get(userId) ?? 0;
    // If scheduled time is stale (past or > 1s ahead), reset to now + 50ms buffer
    const startTime = (scheduled > now && scheduled < now + 1.0) ? scheduled : now + 0.05;
    source.start(startTime);
    this.playbackScheduleTime.set(userId, startTime + audioBuffer.duration);
  }

  addRemoteStream(userId: string, stream: MediaStream) {
    this.remoteStreams.set(userId, stream);

    let audioElement = this.remoteAudioElements.get(userId);
    if (!audioElement) {
      audioElement = new Audio();
      audioElement.autoplay = true;
      audioElement.style.display = 'none';
      document.body.appendChild(audioElement);
      this.remoteAudioElements.set(userId, audioElement);
    }

    audioElement.srcObject = stream;
  }

  removeRemoteStream(userId: string) {
    this.remoteStreams.delete(userId);
    this.playbackScheduleTime.delete(userId);

    const audioElement = this.remoteAudioElements.get(userId);
    if (audioElement) {
      audioElement.pause();
      audioElement.srcObject = null;
      audioElement.remove();
      this.remoteAudioElements.delete(userId);
    }
  }

  setRemoteVolume(userId: string, volume: number) {
    const audioElement = this.remoteAudioElements.get(userId);
    if (audioElement) {
      audioElement.volume = Math.max(0, Math.min(1, volume));
    }
  }

  getLocalVolume(): number {
    if (!this.scriptProcessor || !this.audioContext) return 0;
    const analyser = this.audioContext.createAnalyser();
    this.scriptProcessor.connect(analyser);
    const dataArray = new Uint8Array(analyser.frequencyBinCount);
    analyser.getByteFrequencyData(dataArray);
    return dataArray.reduce((a, b) => a + b) / dataArray.length / 255;
  }

  isAudioContextReady(): boolean {
    return this.audioContext !== null && this.audioContext.state === 'running';
  }

  resumeAudioContext() {
    if (this.audioContext && this.audioContext.state === 'suspended') {
      this.audioContext.resume();
    }
    if (this.playbackContext && this.playbackContext.state === 'suspended') {
      this.playbackContext.resume();
    }
  }

  private float32ToBase64(float32Array: Float32Array): string {
    const bytes = new Uint8Array(float32Array.buffer);
    let binary = '';
    for (let i = 0; i < bytes.byteLength; i++) {
      binary += String.fromCharCode(bytes[i]);
    }
    return btoa(binary);
  }
}

export const audioService = new AudioService();
