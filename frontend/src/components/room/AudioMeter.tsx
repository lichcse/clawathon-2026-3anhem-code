import { useEffect, useRef } from 'react'

interface AudioMeterProps {
  stream: MediaStream | null
  className?: string
}

export function AudioMeter({ stream, className = '' }: AudioMeterProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const animRef = useRef<number>(0)
  const analyserRef = useRef<AnalyserNode | null>(null)
  const ctxRef = useRef<AudioContext | null>(null)

  useEffect(() => {
    if (!stream) return

    const audioCtx = new AudioContext()
    const analyser = audioCtx.createAnalyser()
    analyser.fftSize = 256
    const source = audioCtx.createMediaStreamSource(stream)
    source.connect(analyser)

    ctxRef.current = audioCtx
    analyserRef.current = analyser

    const data = new Uint8Array(analyser.frequencyBinCount)
    const canvas = canvasRef.current
    if (!canvas) return
    const ctx = canvas.getContext('2d')!

    const draw = () => {
      animRef.current = requestAnimationFrame(draw)
      analyser.getByteFrequencyData(data)

      const avg = data.reduce((a, b) => a + b, 0) / data.length
      const level = Math.min(avg / 128, 1)

      ctx.clearRect(0, 0, canvas.width, canvas.height)
      const bars = 12
      const barW = canvas.width / bars - 2
      for (let i = 0; i < bars; i++) {
        const filled = (i + 1) / bars <= level
        ctx.fillStyle = filled ? '#22c55e' : '#374151'
        ctx.fillRect(i * (barW + 2), 0, barW, canvas.height)
      }
    }
    draw()

    return () => {
      cancelAnimationFrame(animRef.current)
      audioCtx.close()
    }
  }, [stream])

  return (
    <canvas
      ref={canvasRef}
      width={80}
      height={12}
      className={className}
      style={{ imageRendering: 'pixelated' }}
    />
  )
}
