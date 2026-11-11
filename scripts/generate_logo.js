const fs = require('fs');
const path = require('path');
const { Resvg } = require('B:/workgit/agentspec/node_modules/@resvg/resvg-js');

function buildLogoSvg() {
  return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024" width="1024" height="1024">
  <defs>
    <clipPath id="squircle-clip">
      <rect x="24" y="24" width="976" height="976" rx="220" />
    </clipPath>

    <linearGradient id="cyan-glow" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#00f5ff"/>
      <stop offset="100%" stop-color="#0284c7"/>
    </linearGradient>

    <linearGradient id="amber-accent" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#fbbf24"/>
      <stop offset="100%" stop-color="#ea580c"/>
    </linearGradient>

    <linearGradient id="slate-light" x1="0%" y1="0%" x2="0%" y2="100%">
      <stop offset="0%" stop-color="#475569"/>
      <stop offset="100%" stop-color="#334155"/>
    </linearGradient>

    <linearGradient id="slate-mid" x1="0%" y1="0%" x2="0%" y2="100%">
      <stop offset="0%" stop-color="#334155"/>
      <stop offset="100%" stop-color="#1e293b"/>
    </linearGradient>

    <linearGradient id="slate-dark" x1="0%" y1="0%" x2="0%" y2="100%">
      <stop offset="0%" stop-color="#1e293b"/>
      <stop offset="100%" stop-color="#0b0f19"/>
    </linearGradient>

    <filter id="subtle-shadow" x="-20%" y="-20%" width="140%" height="140%">
      <feDropShadow dx="0" dy="16" stdDeviation="24" flood-color="#000000" flood-opacity="0.18" />
    </filter>
  </defs>

  <!-- Luxury White Squircle Container -->
  <rect x="24" y="24" width="976" height="976" rx="220" fill="#ffffff" stroke="#e2e8f0" stroke-width="6" />

  <g clip-path="url(#squircle-clip)">
    <g transform="translate(512, 512)" filter="url(#subtle-shadow)">

      <!-- Hexagonal Architectural Gateway Frame -->
      <polygon points="
        0,-410
        355,-205
        355,205
        0,410
        -355,205
        -355,-205
      " fill="none" stroke="#0f172a" stroke-width="36" stroke-linejoin="round" />

      <!-- Inner Dashed Cyan Hexagon -->
      <polygon points="
        0,-375
        324,-187
        324,187
        0,375
        -324,187
        -324,-187
      " fill="none" stroke="#00f5ff" stroke-width="4" opacity="0.45" stroke-dasharray="16, 12" />

      <!-- 100% Watertight Solid Base Silhouette -->
      <path d="
        M 0 -220
        L 45 -220
        L 135 -335
        L 155 -150
        L 205 -50
        L 220 50
        L 180 150
        L 240 270
        L 150 345
        L 0 370
        L -150 345
        L -240 270
        L -180 150
        L -220 50
        L -205 -50
        L -155 -150
        L -135 -335
        L -45 -220
        Z
      " fill="#0b0f19" />

      <!-- EARS -->
      <!-- Left Ear (Facing viewer: left side x < 0) -->
      <polygon points="-45,-220 -135,-335 -155,-150 -95,-130" fill="#1e293b" />
      <polygon points="-55,-215 -125,-320 -95,-140" fill="#0f172a" />
      <polygon points="-75,-200 -120,-305 -90,-150" fill="url(#cyan-glow)" opacity="0.9" />

      <!-- Right Ear (Facing viewer: right side x > 0) -->
      <polygon points="45,-220 135,-335 155,-150 95,-130" fill="#334155" />
      <polygon points="55,-215 125,-320 95,-140" fill="#1e293b" />
      <polygon points="75,-200 120,-305 90,-150" fill="url(#cyan-glow)" opacity="0.95" />

      <!-- CRANIUM & BROW -->
      <polygon points="0,-220 -45,-220 -35,-140 0,-120" fill="#1e293b" />
      <polygon points="0,-220 45,-220 35,-140 0,-120" fill="#475569" />

      <polygon points="-45,-220 -95,-130 -35,-140" fill="#0f172a" />
      <polygon points="45,-220 95,-130 35,-140" fill="#64748b" />

      <!-- Upper Brow Plates -->
      <polygon points="0,-120 -35,-140 -85,-85 -25,-70 0,-80" fill="#1e293b" />
      <polygon points="0,-120 35,-140 85,-85 25,-70 0,-80" fill="#475569" />

      <!-- CHEEKS -->
      <polygon points="-95,-130 -155,-150 -205,-50 -135,-30 -85,-85" fill="#0b0f19" />
      <polygon points="95,-130 155,-150 205,-50 135,-30 85,-85" fill="#1e293b" />

      <polygon points="-205,-50 -220,50 -140,55 -135,-30" fill="#0f172a" />
      <polygon points="205,-50 220,50 140,55 135,-30" fill="#334155" />

      <!-- NOSE BRIDGE (Slate / Titanium Facets with High Contrast) -->
      <polygon points="0,-80 -25,-70 -20,20 0,35" fill="#334155" />
      <polygon points="0,-80 25,-70 20,20 0,35" fill="#64748b" />

      <!-- Lateral Snout Plates -->
      <polygon points="-25,-70 -85,-85 -65,-30 -60,30 -20,20" fill="#1e293b" />
      <polygon points="25,-70 85,-85 65,-30 60,30 20,20" fill="#475569" />

      <!-- ======================================================== -->
      <!-- PREDATOR EYES (Sharp Geometric Almond Optics)            -->
      <!-- ======================================================== -->
      <!-- Left Eye Dark Bezel -->
      <polygon points="-30,-68 -82,-82 -62,-44 -26,-56" fill="#050810" />
      <!-- Left Eye Iris (Electric Cyan) -->
      <polygon points="-34,-66 -76,-78 -60,-48 -30,-58" fill="url(#cyan-glow)" />
      <!-- Left Pupil & Slit Core -->
      <polygon points="-48,-68 -56,-68 -52,-54 -46,-56" fill="#ffffff" />

      <!-- Right Eye Dark Bezel -->
      <polygon points="30,-68 82,-82 62,-44 26,-56" fill="#050810" />
      <!-- Right Eye Iris (Electric Cyan) -->
      <polygon points="34,-66 76,-78 60,-48 30,-58" fill="url(#cyan-glow)" />
      <!-- Right Pupil & Slit Core -->
      <polygon points="48,-68 56,-68 52,-54 46,-56" fill="#ffffff" />

      <!-- Subtle Cyan Eye Reflection Nodes -->
      <circle cx="-52" cy="-60" r="3.5" fill="#ffffff" />
      <circle cx="52" cy="-60" r="3.5" fill="#ffffff" />

      <!-- NOSE PAD (Hexagonal Obsidian Mark) -->
      <polygon points="0,35 -24,23 -18,52 0,60" fill="#090d16" />
      <polygon points="0,35 24,23 18,52 0,60" fill="#1e293b" />

      <!-- MUZZLE FLAPS & WHISKER PADS -->
      <polygon points="0,60 -18,52 -55,70 -25,100 0,110" fill="#0f172a" />
      <polygon points="0,60 18,52 55,70 25,100 0,110" fill="#334155" />

      <polygon points="-55,70 -60,30 -140,55 -95,120 -25,100" fill="#0b0f19" />
      <polygon points="55,70 60,30 140,55 95,120 25,100" fill="#1e293b" />

      <!-- POINTED CHIN APEX -->
      <polygon points="0,110 -25,100 -16,140 0,150" fill="#090d16" />
      <polygon points="0,110 25,100 16,140 0,150" fill="#1e293b" />

      <!-- ======================================================== -->
      <!-- GUARDIAN GORGET / DAG COLLAR & CHEST PLATES             -->
      <!-- ======================================================== -->

      <!-- Throat Apex -->
      <polygon points="0,150 -16,140 -60,170 0,205" fill="#1e293b" />
      <polygon points="0,150 16,140 60,170 0,205" fill="#334155" />

      <!-- Lateral Throat Plates -->
      <polygon points="-16,140 -95,120 -145,170 -60,170" fill="#090d16" />
      <polygon points="16,140 95,120 145,170 60,170" fill="#1e293b" />

      <!-- Flank Armor Plates -->
      <polygon points="-140,55 -180,150 -145,170 -95,120" fill="#0b0f19" />
      <polygon points="140,55 180,150 145,170 95,120" fill="#334155" />

      <!-- PRIMARY DAG CORE COLLAR (Obsidian with Cyan/Amber Node Matrix) -->
      <polygon points="0,205 -60,170 -100,235 0,280" fill="url(#slate-dark)" />
      <polygon points="0,205 60,170 100,235 0,280" fill="url(#slate-mid)" />

      <!-- Interlocking DAG Flow Chevron -->
      <polygon points="0,220 -40,195 -25,240 0,260" fill="url(#cyan-glow)" />
      <polygon points="0,220 40,195 25,240 0,260" fill="url(#cyan-glow)" />

      <!-- Network DAG Visual Connectors (Collar Flow Vectors) -->
      <line x1="-32" y1="210" x2="0" y2="235" stroke="#ffffff" stroke-width="3" stroke-linecap="round" />
      <line x1="32" y1="210" x2="0" y2="235" stroke="#ffffff" stroke-width="3" stroke-linecap="round" />
      <line x1="0" y1="235" x2="0" y2="265" stroke="#ffffff" stroke-width="3" stroke-linecap="round" />

      <!-- DAG Core Nodes -->
      <circle cx="0" cy="235" r="7" fill="#ffffff" />
      <circle cx="-32" cy="210" r="5" fill="#00f5ff" />
      <circle cx="32" cy="210" r="5" fill="#00f5ff" />
      <circle cx="0" cy="265" r="5" fill="#00f5ff" />

      <!-- Amber Security Alert Status Nodes on Armor Flanks -->
      <polygon points="-100,235 -130,230 -115,255" fill="url(#amber-accent)" />
      <polygon points="100,235 130,230 115,255" fill="url(#amber-accent)" />

      <!-- Low-Poly Breastplate -->
      <polygon points="0,280 -100,235 -170,260 -130,325 0,360" fill="#090d16" />
      <polygon points="0,280 100,235 170,260 130,325 0,360" fill="#1e293b" />

      <!-- Shoulder Pauldrons -->
      <polygon points="-145,170 -180,150 -240,270 -170,260" fill="#0f172a" />
      <polygon points="145,170 180,150 240,270 170,260" fill="#334155" />

    </g>
  </g>
</svg>`;
}

async function renderLogo(outputDir) {
  const svg = buildLogoSvg();
  const svgPath = path.join(outputDir, 'logo.svg');
  const pngPath = path.join(outputDir, 'logo.png');

  fs.writeFileSync(svgPath, svg, 'utf8');
  const resvg = new Resvg(svg, { fitTo: { mode: 'width', value: 1024 } });
  fs.writeFileSync(pngPath, resvg.render().asPng());
  console.log(`✓ Successfully rendered refined logo.svg and logo.png at 1024x1024 in ${outputDir}`);
}

const targetDir = path.resolve('B:/workgit/n8nflow-guard/docs/images');
renderLogo(targetDir).catch(err => {
  console.error('Error rendering logo:', err);
  process.exit(1);
});
