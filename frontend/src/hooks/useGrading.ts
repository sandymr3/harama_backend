export const useGrading = (submissionId: string) => {
    return {
        grading: {
            score: 85,
            confidence: 0.9,
            aiReasoning: "Good job. The diagram of the plant cell is mostly accurate, though the mitochondria labels are a bit blurry.",
            answer: {
                text: "The cell membrane protects the cell and controls what goes in and out. The nucleus contains DNA.",
                diagrams: [
                    "https://placehold.co/400x300?text=Plant+Cell+Diagram",
                    "https://placehold.co/400x300?text=Nucleus+Detail"
                ]
            }
        },
        loading: false,
        applyOverride: (score: number) => console.log(score)
    }
}
