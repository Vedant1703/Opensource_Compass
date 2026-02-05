"use client";

import { useEffect, useRef } from "react";
import { Button } from "@/components/ui/button";
import { CloudLightning, Github, Lightbulb } from "lucide-react";
import gsap from "gsap";

export default function HeroSection() {
  const containerRef = useRef<HTMLDivElement>(null);

  function SplitText({ text }: { text: string }) {
    return (
      <>
        {text.split("").map((char, i) => (
          <span key={i} className="hero-letter inline-block" aria-hidden>
            {char === " " ? "\u00A0" : char}
          </span>
        ))}
      </>
    );
  }

  useEffect(() => {
    const ctx = gsap.context(() => {
      const tl = gsap.timeline();

      tl.from(".hero-letter", {
        y: -120,
        opacity: 0,
        filter: "blur(10px)",
        duration: 0.9,
        stagger: 0.03,
        ease: "back.out(1.8)",
      })
        .from(
          ".hero-sub",
          {
            y: 20,
            opacity: 0,
            duration: 0.6,
            ease: "power3.out",
          },
          "-=0.4",
        )
        .from(
          ".hero-cta",
          {
            y: 16,
            opacity: 0,
            scale: 0.96,
            duration: 0.5,
            ease: "power3.out",
          },
          "-=0.3",
        )
        .from(
          ".hero-foot",
          {
            opacity: 0,
            duration: 0.4,
          },
          "-=0.2",
        );
    }, containerRef);

    return () => ctx.revert();
  }, []);

  return (
    <section
      ref={containerRef}
      className="relative min-h-[90vh] flex items-center justify-center bg-[#0d1117] text-white overflow-hidden"
    >
      {/* Subtle background gradient */}
      <div className="absolute inset-0 bg-gradient-to-b from-[#0d1117] via-[#0d1117] to-[#010409]" />

      <div className="relative z-10 max-w-6xl w-full px-6 py-20 grid grid-cols-1 md:grid-cols-2 gap-14 items-center">
        {/* LEFT CONTENT */}
        <div className="space-y-6 text-center md:text-left">
<h1 className="text-4xl md:text-5xl xl:text-6xl font-semibold tracking-tight leading-none ">
  <span className="block py-4 overflow-hidden">
    <SplitText text="Start Open Source" />
  </span>

  <span className="block overflow-hidden -mt-2">
    <SplitText text="the " />
    <span className="text-[#2f81f7] glow-text inline-block">
      <SplitText text="Right Way" />
    </span>
  </span>
</h1>


          <p className="hero-sub text-base md:text-lg text-[#8b949e] max-w-xl mx-auto md:mx-0">
            Discover beginner-friendly open source projects and issues tailored
            to your tech stack — guided by AI, backed by GitHub data.
          </p>

          <div className="hero-cta flex flex-col sm:flex-row gap-4 justify-center md:justify-start">
            <Button
              size="lg"
              className="bg-[#238636] hover:bg-[#2ea043] text-white px-6 py-5 text-base font-medium"
           onClick={() => { window.location.href = "/onboarding"; }}
>
              Get Started
            </Button>
          </div>

          <p className="hero-foot text-xs text-[#6e7681]">
            Read-only GitHub access • No spam • Open-source friendly
          </p>
        </div>

        {/* RIGHT VISUAL (Optional) */}
        <div className="hidden md:block">
          <div className="hero-animate relative rounded-xl border border-[#30363d] bg-[#161b22] p-6 shadow-lg">
            <div className="space-y-3">
              <div className="h-3 w-1/3 rounded bg-[#30363d]" />
              <div className="h-3 w-2/3 rounded bg-[#30363d]" />
              <div className="h-3 w-1/2 rounded bg-[#30363d]" />

              <div className="mt-6 space-y-2">
                <div className="h-10 rounded bg-[#0d1117] border border-[#30363d]" />
                <div className="h-10 rounded bg-[#0d1117] border border-[#30363d]" />
                <div className="h-10 rounded bg-[#0d1117] border border-[#30363d]" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
