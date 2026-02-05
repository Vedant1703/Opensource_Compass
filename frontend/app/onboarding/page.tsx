"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import gsap from "gsap";
import AuthCard from "./components/authcard";
import ProgressStepper from "./components/progress";
import TechSection from "./components/techsection";
import StickyActionBar from "./components/stickyactionbar";
import FrameworkSection from "./components/frameworksection";
import DomainSection from "./components/domainsection";


export default function OnboardingPage() {
    const router = useRouter();
    const [frameworks, setFrameworks] = useState<string[]>([]);
    const [domains, setDomains] = useState<string[]>([]);
    const pageRef = useRef<HTMLDivElement>(null);
    const [selectedLanguages, setSelectedLanguages] = useState<string[]>([]);
    
    // Auth state management
    const [isAuth, setIsAuth] = useState(false);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        // Check for auth token
        const token = localStorage.getItem("authToken");
        setIsAuth(!!token);
        setIsLoading(false);
    }, []);

    const currentStep = isAuth ? 2 : 1;

    useEffect(() => {
        if (isLoading) return; // Don't animate until we know the state
        
        const ctx = gsap.context(() => {
            gsap.from(".onboard-animate", {
                opacity: 0,
                y: 24,
                duration: 0.8,
                stagger: 0.15,
                ease: "power3.out",
            });
        }, pageRef);

        return () => ctx.revert();
    }, [isLoading, isAuth]); // Re-run animation when state changes

    const handleContinue = () => {
        router.push("/discover");
    };

    if (isLoading) return null; // Or a loading spinner

    return (
        <div
            ref={pageRef}
            className="min-h-screen bg-[#0d1117] text-white flex justify-center pb-24"
        >
            <div className="w-full max-w-3xl px-6 py-16 space-y-12">

                {/* Step 1: Auth Card - Only show if NOT authenticated */}
                {!isAuth && (
                    <div className="onboard-animate">
                        <AuthCard />
                    </div>
                )}

                <div className="onboard-animate">
                    <ProgressStepper currentStep={currentStep} />
                </div>

                {/* Step 2: Tech Stack - Only show if authenticated */}
                {isAuth && (
                    <>
                        <div className="onboard-animate">
                            <TechSection
                                title="Programming Languages"
                                subtitle="Select at least one language you are comfortable with"
                                options={[
                                    "JavaScript",
                                    "TypeScript",
                                    "Python",
                                    "Java",
                                    "Go",
                                    "Rust",
                                    "C++",
                                ]}
                                required
                                selected={selectedLanguages}
                                setSelected={setSelectedLanguages}
                            />
                            <FrameworkSection
                                selected={frameworks}
                                setSelected={setFrameworks}
                            />
                        </div>

                        <div className="onboard-animate">
                            <DomainSection
                                selected={domains}
                                setSelected={setDomains}
                            />
                        </div>
                    </>
                )}

                {/* Sticky Action Bar - Only show if authenticated (Step 2) */}
                {isAuth && (
                    <StickyActionBar 
                        isValid={selectedLanguages.length > 0} 
                        onContinue={handleContinue}
                    />
                )}
            </div>
        </div>
    );
}
