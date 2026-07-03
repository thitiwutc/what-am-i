<script setup lang="ts">
import { CommandType } from '@/types/command-type'
import { getEnv } from '@/env'
import { EventType } from '@/types/event-type'
import type { Room } from '@/types/room'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import type { Player } from '@/types/player'
import { Role } from '@/types/role-type'
import { RoomState } from '@/types/room-state'
import type { PlayerState } from '@/types/player-state'

const env = getEnv()
const ws = new WebSocket(new URL('ws', env.apiUrl))
const curUrl = window.location.href
const room = ref<Room>({
  id: '',
  state: RoomState.Lobby,
  players: [],
})
const pid = ref('')
const player = computed((): Player => {
  const zeroValuePlayer: Player = {
    id: '',
    name: '',
    role: Role.PARTICIPANT,
    isSecretWordSubmitted: false,
  }
  if (!pid.value) {
    return zeroValuePlayer
  }

  const found: Player | undefined = room.value.players.find((p) => p.id === pid.value)
  if (!found) {
    return zeroValuePlayer
  }

  console.log({ found })
  return found
})
const secretWord = ref('')
const route = useRoute()

onMounted(async () => {
  // Connection opened
  ws.onopen = function () {
    // Join room
    ws.send(
      JSON.stringify({
        type: CommandType.JoinRoom,
        rid: route.params.room_id as string,
        pname: sessionStorage.getItem('pname') ?? '',
      }),
    )
  }

  // Listen for messages
  ws.onmessage = function (event) {
    console.log('Received websocket message:', event.data)
    const serverEvent = JSON.parse(event.data)
    console.debug({ serverEvent })
    switch (serverEvent.type) {
      case EventType.Error:
        console.error(serverEvent.error)
        break
      case EventType.RoomState:
        // Update room state
        room.value = {
          id: serverEvent.room_state.id,
          state: serverEvent.room_state.state,
          players: serverEvent.room_state.players.map((p: PlayerState) => ({
            id: p.id,
            name: p.name,
            role: p.role,
            isSecretWordSubmitted: p.isws,
          })),
        }
        break
      case EventType.JoinRoomResult:
        pid.value = serverEvent.player.id
        break
    }
  }

  // Handle errors
  ws.onerror = function (event) {
    console.error('Webws error:', event)
  }

  // Handle disconnection
  ws.onclose = function (event) {
    if (event.wasClean) {
      console.log(`Closed cleanly, code=${event.code}, reason=${event.reason}`)
    } else {
      console.log('Connection died')
    }
  }
})

onUnmounted(() => {
  console.log('Closing websocket connection')
  ws.close()
})

function startGame() {
  console.debug('Called startGame')
  ws.send(
    JSON.stringify({
      type: CommandType.StartGame,
      rid: room.value.id,
    }),
  )
}

function submitSecretWord(ev: Event) {
  if (!secretWord.value) {
    return
  }
  if (ev instanceof KeyboardEvent && ev.key != 'Enter') {
    return
  }

  ws.send(
    JSON.stringify({
      type: CommandType.SubmitSecretWord,
      rid: room.value.id,
      pid: player.value.id,
      word: secretWord.value,
    }),
  )

  console.log(`Secret word "${secretWord.value}" submitted`)
}
</script>

<template>
  <div class="container d-flex justify-content-center align-items-center vh-100">
    <div v-if="room.state === RoomState.Lobby" class="row w-100 justify-content-center">
      <div
        class="col-12 col-sm-10 col-md-8 col-lg-6 col-xl-4 shadow p-3 rounded border"
        style="background: white"
      >
        <div class="row gy-2">
          <h2 class="col-12">Room ID: {{ room.id }}</h2>
          <div class="col-12">
            <div class="row">
              <div class="col-auto d-flex align-items-center">Invite link:</div>
              <div class="col">
                <input class="form-control" type="text" readonly :value="curUrl" />
              </div>
            </div>
          </div>
          <div class="col-12">
            <div class="row">
              <div class="col-12">
                <b>{{ room.players.length > 1 ? 'Players' : 'Players' }}</b>
              </div>
              <div class="col-12" v-for="player in room.players" :key="player.id">
                <div v-if="player.role === Role.HOST">
                  {{ player.name }}
                  <i class="bi bi-yin-yang"></i>
                </div>
                <div v-else>
                  {{ player.name }}
                </div>
              </div>
              <div class="div d-flex justify-content-end">
                <button
                  class="btn btn-primary"
                  :disabled="room.players.length < 2 || player.role !== Role.HOST"
                  @click="startGame"
                >
                  Start
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Playing state -->
    <div class="row justify-content-center" v-else-if="room.state === RoomState.Thinking">
      <div class="col rounded border" style="background: white">
        <h3 v-if="player.isSecretWordSubmitted">
          Waiting for other players to submit their secret words
        </h3>
        <div v-else class="container p-3">
          <div class="row">
            <div class="col">
              <input
                class="form-control"
                type="text"
                v-model="secretWord"
                @keydown="submitSecretWord"
                :disabled="player.isSecretWordSubmitted"
                placeholder="Secret word"
              />
            </div>
          </div>
          <div class="row">
            <div class="mt-1 d-flex justify-content-end">
              <button class="btn btn-primary" @click="submitSecretWord">Submit</button>
            </div>
          </div>
        </div>
        <div class="mt-2">
          <ul>
            <li v-for="player in room.players" :key="player.id">
              {{ player.name }} <i v-show="player.isSecretWordSubmitted" class="bi bi-check"></i>
            </li>
          </ul>
        </div>
      </div>
    </div>

    <div v-else-if="room.state === RoomState.Guessing">Guessing</div>
    <div v-else-if="room.state === RoomState.Finished">Finished</div>
    <div v-else>
      <h2>Invalid room state: {{ room.state }}</h2>
    </div>
  </div>
</template>
