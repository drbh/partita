<script lang="ts">
  import { Canvas } from "@threlte/core";
  import { Grid } from "@threlte/extras";
  import { onMount } from "svelte";
  import Scene from "./scene.svelte";

  const DIRECTIONS = {
    FRONT: parseFloat((2 * Math.PI).toFixed(4)),
    LEFT: parseFloat((-Math.PI / 2).toFixed(4)),
    BACK: parseFloat((2 * Math.PI + Math.PI).toFixed(4)),
    RIGHT: parseFloat((2 * Math.PI + Math.PI / 2).toFixed(4)),
  };

  let gameStarted = false;

  let player = {
    y: 0,
    rotation: 0,
    ghostRotation: DIRECTIONS.FRONT,
    position: [0, 0, 0],
    isLookingForGame: false,
    gameKey: "drbh",
    showModal: true,
    inGameModal: false,
    step: 0,
    id: "",
    positionPath: [],
  };

  let opponent = {
    ghostRotation: DIRECTIONS.LEFT,
    position: [7, player.y, 0],
    positionPath: [],
  };

  let socket: WebSocket;

  onMount(() => {
    setupWebSocket();
    setupKeydownEvent();
  });

  function setupWebSocket() {
    let host =
      window.location.hostname === "localhost"
        ? "ws://localhost:3000"
        : "wss://partita.fly.dev";

    console.log("connecting to", host);

    socket = new WebSocket(`${host}/ws/game`);

    socket.addEventListener("message", handleSocketMessage);
  }

  function handleSocketMessage(event: MessageEvent) {
    try {
      const data = JSON.parse(event.data);

      if (data.command && data.command === "matchFound") {
        player.gameKey = data.gameKey;
        player.isLookingForGame = false;
      }

      if (data.command && data.command === "playerNameSet") {
        // validation successful and name is set
        nextStep();
        return;
      }
      if (data.command && data.command === "playerCollision") {
        let newLog = {
          username: data.name,
          time: /* random string */ Math.random().toString(36).substring(10),
          content: `collided with ${data.with}`,
        };

        // add to activity log
        activityLog = [...activityLog.slice(-2), newLog];
        return;
      }

      // if no player data, return
      if (!data[player.gameKey]) {
        return;
      }

      if (!data[player.gameKey].Players) {
        return;
      }

      // if player data is empty, return
      if (Object.keys(data[player.gameKey].Players).length === 0) {
        return;
      }

      // if player is not in game, return
      if (!data[player.gameKey].Players[player.id]) {
        return;
      }

      updatePlayerPosition(data);
      updateOpponentPosition(data);
    } catch (error) {
      console.log(error);
    }
  }

  function updatePlayerPosition(data: any) {
    const pos = data[player.gameKey].Players[player.id];

    player.position[0] = pos.X;
    player.position[2] = pos.Z;

    player.positionPath = pos.PathPoints.map((p: any) => [p.X, p.Y, p.Z]);

    player.ghostRotation = pos.Rotation;
  }

  function updateOpponentPosition(data: any) {
    const players = Object.keys(data[player.gameKey].Players);
    if (players.length > 1) {
      const otherPlayerId = players.filter((p) => p !== player.id)[0];
      const otherPlayer = data[player.gameKey].Players[otherPlayerId];
      opponent.position[0] = otherPlayer.X;
      opponent.position[2] = otherPlayer.Z;
      opponent.ghostRotation = otherPlayer.Rotation;
      opponent.positionPath = otherPlayer.PathPoints.map((p: any) => [
        p.X,
        p.Y,
        p.Z,
      ]);
    }
  }

  function setupKeydownEvent() {
    window.addEventListener("keydown", (event) => {
      let targetRotation = 0;
      if (event.key === "ArrowLeft") {
        targetRotation = DIRECTIONS.LEFT;
      } else if (event.key === "ArrowRight") {
        targetRotation = DIRECTIONS.RIGHT;
      } else if (event.key === "ArrowUp") {
        targetRotation = DIRECTIONS.BACK;
      } else if (event.key === "ArrowDown") {
        targetRotation = DIRECTIONS.FRONT;
      }
      if (targetRotation === 0) {
        return;
      }
      socket.send("rotate:" + targetRotation);
    });
  }

  function startGame() {
    socket.send(`startGame:${player.gameKey}:100`);
    player.showModal = false;
    gameStarted = true;
  }

  function joinGame() {
    socket.send(`joinGame:${player.gameKey}`);
    player.showModal = false;
    gameStarted = true;
  }

  function leaveGame() {
    socket.send(`leaveGame:${player.gameKey}`);
    player.showModal = false;
    gameStarted = false;
  }

  function findGame() {
    socket.send(`findGame`);
    player.isLookingForGame = true;
  }

  function openModal() {
    player.showModal = true;
  }

  function closeModal() {
    player.showModal = false;
  }

  function closeInGameModal() {
    player.inGameModal = false;
  }

  function nextStep() {
    player.step++;
  }

  function updateUsername() {
    localStorage.setItem("username", player.id);
  }

  let errorMessage = "";

  function validateUsername() {
    // Dummy validation, replace with actual validation logic.
    if (player.id.length < 3) {
      errorMessage = "Username should be at least 3 characters long!";
    } else if (false) {
      errorMessage = "Username is already taken!";
    } else {
      socket.send(`setPlayerName:${player.id}`);
    }
  }

  let activityLog: any[] = [];
</script>

<div class="settings-button">
  <button class="ball" on:click={openModal}>SETTINGS</button>
</div>

<!-- Game Activity Log -->
<div class="activity-log">
  <div class="activity-log-header">
    <h4>Activity Log</h4>
  </div>
  <div class="activity-log-content">
    {#each activityLog as log (log.time)}
      <div class="activity-log-item">
        <div class="activity-log-item-header">
          <div class="activity-log-item-header-icon">👻</div>
          <div class="activity-log-item-header-text">
            <div class="activity-log-item-header-text-username">
              <span>{log.username}</span>
            </div>
            <div class="activity-log-item-header-text-time">
              <span>{log.time}</span>
            </div>
          </div>
        </div>
        <div class="activity-log-item-content">
          <span>{log.content}</span>
        </div>
      </div>
    {/each}
  </div>
</div>

{#if player.showModal}
  <div class="modal">
    <div class="modal-content">
      <button class="close" on:click={closeModal}>
        <span>&times;</span>
      </button>

      {#if gameStarted}
        <button class="ball" on:click={leaveGame}>Leave Game</button>
      {:else}
        <div class="modal-header">
          <h2>Game Settings</h2>
        </div>

        {#if player.step === 0}
          <div class="welcome-screen">
            <h1>Welcome to Partita!</h1>
            <p>
              Immerse yourself in a world of adventure and fun. Let's get
              started!
            </p>
            <button class="ball" on:click={nextStep}>Begin Journey</button>
          </div>
        {/if}

        {#if player.step === 1}
          <div class="username-input">
            <label for="username">Choose your username to get started:</label>
            <input
              type="text"
              id="username"
              name="username"
              bind:value={player.id}
              autofocus
            />
            <button class="ball" on:click={validateUsername}>Continue</button>
          </div>
        {/if}

        {#if player.step === 2}
          <div class="main-menu">
            <button
              class="ball"
              style="margin-bottom: 20px;"
              on:click={startGame}>Start a New Game</button
            >
            <button
              class="ball"
              style="margin-bottom: 20px;"
              on:click={findGame}>Find a Match</button
            >
            <div class="join-game">
              <label for="gameKey">Have a game key? Join a game:</label>
              <div class="join-game-input">
                <input
                  type="text"
                  id="gameKey"
                  name="gameKey"
                  bind:value={player.gameKey}
                  placeholder="Enter Game Key"
                />
                <button class="ball" on:click={joinGame}>Join</button>
              </div>
            </div>
          </div>
        {/if}

        {#if errorMessage}
          <div class="error-message">{errorMessage}</div>
        {/if}

        {#if player.isLookingForGame}
          <div class="loading-message">Finding games...</div>
        {/if}
      {/if}
    </div>
  </div>
{/if}

<!-- modal for in game settings -->
{#if player.inGameModal}
  <div class="modal">
    <div class="modal-content">
      <button class="close" on:click={closeInGameModal}>
        <span>&times;</span>
      </button>

      <div class="modal-header">
        <h2>In-Game Settings</h2>
      </div>

      <button class="ball" on:click={leaveGame}>Leave Game</button>
    </div>
  </div>
{/if}

<Canvas>
  <Grid>
    <Scene
      bind:rotation={player.rotation}
      bind:ghostRotation={player.ghostRotation}
      bind:position={player.position}
      bind:ghostRotation2={opponent.ghostRotation}
      bind:position2={opponent.position}
      bind:playerPostionPath={player.positionPath}
      bind:player2PostionPath={opponent.positionPath}
    />
  </Grid>
</Canvas>

<style>
  * {
    font-family: "Poppins", sans-serif;
    color: #888;
  }

  .settings-button {
    position: absolute;
    top: 10px;
    right: 30px;
    z-index: 1000;
    display: flex;
    flex-direction: column;
    color: #888;
  }

  .ball {
    padding: 8px 16px;
    background-color: #ff6347;
    color: white;
    font-size: 16px;
    border-radius: 50px;
    border: none;
    box-shadow: 0px 8px 15px rgba(0, 0, 0, 0.1);
    font-weight: bold;
    cursor: pointer;
    z-index: 2;
  }

  .ball:hover {
    background-color: #ff4500;
  }

  .modal {
    position: absolute;
    z-index: 3;
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;
    overflow: auto;
    background-color: rgba(79, 79, 79, 0.4);
  }

  .modal-content {
    background-color: #1d1d1d;
    color: #888;
    margin: 15% auto;
    max-width: 600px;
    border: none;
    border-radius: 10px;
    padding: 30px;
    width: 80%;
    box-shadow: 0px 8px 15px rgba(0, 0, 0, 0.1);
  }

  .modal-header {
    padding-bottom: 20px;
  }

  .main-menu {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .close {
    color: #aaa;
    float: right;
    font-size: 28px;
    font-weight: bold;
  }

  .close:hover,
  .close:focus {
    color: #333;
    text-decoration: none;
    cursor: pointer;
  }

  input[type="text"] {
    border: none;
    border-radius: 50px;
    background: #f2f2f2;
    padding: 10px 20px;
    color: #333;
    font-size: 16px;
    font-weight: bold;
    margin-bottom: 20px;
  }

  .activity-log {
    position: absolute;
    top: 10px;
    left: 10px;
    z-index: 1;
    display: flex;
    flex-direction: column;
    color: #888;
    background: #1d1d1d33;
    border-radius: 5px;
    box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.1);
    width: 300px;
    padding: 5px;
  }

  .activity-log-header {
    padding-bottom: 5px;
    font-size: 12px;
    font-weight: normal;
    color: #fff;
  }

  .activity-log-item {
    display: flex;
    flex-direction: column;
    padding: 5px;
    border-bottom: 1px solid #444;
    color: #aaa;
  }

  .activity-log-item-header {
    display: flex;
    flex-direction: row;
    align-items: center;
    padding-bottom: 5px;
  }

  .activity-log-item-header-icon {
    padding-right: 5px;
    font-size: 10px;
  }
</style>
