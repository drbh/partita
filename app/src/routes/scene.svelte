<script lang="ts">
  import { T } from "@threlte/core";
  import {
    MeshLineGeometry,
    MeshLineMaterial,
    useGltf
  } from "@threlte/extras";
  import * as Three from "three";

  const ghostScale = 0.25;
  
  let y = 0;
  let boundary = 8;
  let pathHeight = 0;

  export let rotation = 0;
  export let ghostRotation = 0;
  export let position = [0, y, 0];
  export let ghostRotation2 = 0;
  export let position2 = [7, y, 0];
  export let playerPostionPath: number[][] = [];
  export let player2PostionPath: number[][] = [];

  // Helper function to generate line data
  function generateLineData(
    path: number[][],
    currentPos: number[]
  ): Three.Vector3[] {
    let lineData: Three.Vector3[] = Array.from(
      { length: 100 },
      () => new Three.Vector3(0, 0, 0)
    );
    if (path.length > 2) {
      const offset = lineData.length - path.length;

      // overwrite the all the points with the first point
      lineData.fill(new Three.Vector3(...path[0]), 0, offset);

      for (let i = 0; i < path.length; i++) {
        lineData[i + offset - 1] = new Three.Vector3(...path[i]);
        lineData[i + offset - 1].y = pathHeight;
      }
      // add the current position to the end of the line
      lineData[lineData.length - 1] = new Three.Vector3(...currentPos);
      lineData[lineData.length - 1].y = pathHeight;
    }
    return lineData;
  }

  // Generate line data for player 1
  $: lineData = generateLineData(
    [...playerPostionPath, [position[0], 0.25, position[2]]],
    position
  );

  // Generate line data for player 2
  $: lineData2 = generateLineData(
    [...player2PostionPath, [position2[0], 0.25, position2[2]]],
    position2
  );

  const gltf = useGltf<{
    nodes: {
      Ghost002: THREE.Mesh;
    };
    materials: {
      ["Ghost002.001"]: THREE.MeshStandardMaterial;
    };
  }>("/assets/ghost2.glb");

  $: if ($gltf) {
    const objectA = $gltf.nodes["Ghost002"];
    // @ts-ignore-next-line
    objectA.material.color.set("#00d062");
  }
</script>

<T.PerspectiveCamera
  makeDefault
  position={[position[0] + 0, position[1] + 20, position[2] + 20]}
  fov={36}
  target={position}
  on:create={({ ref }) => {
    ref.lookAt(0, 0, 0);
  }}
>
  <!-- <OrbitControls /> -->
</T.PerspectiveCamera>

<T.AmbientLight color="#9393ac" intensity={10} />
<T.PointLight intensity={10} position={[4, 2, 4]} color="#76aac8" />

{#await useGltf("/assets/ghost.glb") then ghost}
  <T
    is={ghost.scene}
    {position}
    scale={ghostScale}
    rotation.y={ghostRotation}
  />
{/await}

{#await gltf then ghost2}
  <T
    is={ghost2.scene}
    position={position2}
    scale={ghostScale}
    rotation.y={ghostRotation2}
  />
{/await}

<!-- <Grid /> -->

<T.Mesh position={[-boundary, 0, -boundary]} let:ref castShadow>
  <T.SphereGeometry args={[0.1, 32, 32]} />
  <T.MeshStandardMaterial color="white" />
</T.Mesh>

<T.Mesh position={[boundary, 0, -boundary]} let:ref castShadow>
  <T.SphereGeometry args={[0.1, 32, 32]} />
  <T.MeshStandardMaterial color="white" />
</T.Mesh>

<T.Mesh position={[-boundary, 0, boundary]} let:ref castShadow>
  <T.SphereGeometry args={[0.1, 32, 32]} />
  <T.MeshStandardMaterial color="white" />
</T.Mesh>

<T.Mesh position={[boundary, 0, boundary]} let:ref castShadow>
  <T.SphereGeometry args={[0.1, 32, 32]} />
  <T.MeshStandardMaterial color="white" />
</T.Mesh>

<T.Mesh>
  <MeshLineGeometry points={lineData} />
  <MeshLineMaterial width={0.05} color="#b624bb" />
</T.Mesh>

<T.Mesh>
  <MeshLineGeometry points={lineData2} />
  <MeshLineMaterial width={0.05} color="#edcf6b" />
</T.Mesh>

<T.Mesh position={[0, -0.025, 0]} rotation.x={-Math.PI / 2} castShadow>
  <T.PlaneGeometry args={[boundary * 2, boundary * 2, 1, 1]} />
  <T.MeshStandardMaterial color="#323538" />
</T.Mesh>
