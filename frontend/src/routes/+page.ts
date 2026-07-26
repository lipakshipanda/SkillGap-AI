// ApexCharts and the drag-and-drop file upload both touch window/document
// directly, so this route must be rendered client-side only.
export const ssr = false
