export class AudioService {
  private audioContext: AudioContext | null = null;
  private mediaStreamAudioSourceNode: MediaStreamAudioSourceNode | null = null;
  private scriptProcessor: ScriptProcessorNode | null = null;
  private localStream: MediaStream | null = null;
  private remoteAudioElements: Map<string, HTMLAudioElement> = new Map();
  private remoteStreams: Map<string, MediaStream> = new Map();

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

      this.mediaStreamAudioSourceNode.connect(this.scriptProcessor);
      this.scriptProcessor.connect(this.audioContext.destination);

      return this.localStream;
    } catch (error) {
      console.error('Error accessing microphone:', error);
      throw error;
    }
  }

  async stopCapture() {
    if (this.localStream) {
      this.localStream.getTracks().forEach((track) => track.stop());
      this.localStream = null;
    }

    if (this.scriptProcessor) {
      this.scriptProcessor.disconnect();
    }

    if (this.mediaStreamAudioSourceNode) {
      this.mediaStreamAudioSourceNode.disconnect();
    }
  }

  setMicMuted(isMuted: boolean) {
    if (this.localStream) {
      this.localStream.getAudioTracks().forEach((track) => {
        track.enabled = !isMuted;
      });
    }
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
    if (!this.scriptProcessor) return 0;
    const analyser = this.audioContext!.createAnalyser();
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
  }
}

export const audioService = new AudioService();
