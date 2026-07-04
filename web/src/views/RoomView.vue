<script setup lang="ts">
import { CommandType } from '@/types/command-type'
import { getEnv } from '@/env'
import { EventType } from '@/types/event-type'
import type { Room } from '@/types/room'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import type { Player } from '@/types/player'
import { Role } from '@/types/role-type'
import { RoomState } from '@/types/room-state'
import { Answer } from '@/types/answer'
import { Desicion } from '@/types/desicion'

const env = getEnv()
const ws = new WebSocket(new URL('ws', env.apiUrl))
const curUrl = window.location.href
const room = ref<Room>({
  id: '',
  state: RoomState.Lobby,
  players: [],
  active_player: {
    id: '',
    name: '',
    role: Role.PARTICIPANT,
    secret_word: '',
    isws: false,
    questions: [],
    is_word_guessed: false,
  },
  last_question: null,
  last_guess: null,
})
const pid = ref('')
const player = computed((): Player => {
  const zeroValuePlayer: Player = {
    id: '',
    name: '',
    role: Role.PARTICIPANT,
    secret_word: '',
    isws: false,
    questions: [],
    is_word_guessed: false,
  }
  if (!pid.value) {
    return zeroValuePlayer
  }

  const found: Player | undefined = room.value.players.find((p) => p.id === pid.value)
  if (!found) {
    return zeroValuePlayer
  }

  return {
    id: found.id,
    name: found.name,
    role: found.role,
    secret_word: found.secret_word,
    isws: found.isws,
    questions: found.questions,
    is_word_guessed: found.is_word_guessed,
  }
})
const secretWord = ref('')
const inviteLinkCopied = ref(false)
const askGuessText = ref('')
const route = useRoute()
const alertCorrectGuess = ref(false)

watch(room, (room) => {
  if (room.state == RoomState.Finished) {
    console.log('Gate is finished. Closing websocket connection')
    ws.close()
  }
})

watch(player, (newPlayer) => {
  console.debug({ player: newPlayer })
  if (player.value.is_word_guessed && !alertCorrectGuess.value) {
    alertCorrectGuess.value = true
    alert('You guess the word correctly!')
  }
})

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
    console.debug({ websocketMsg: serverEvent })
    switch (serverEvent.type) {
      case EventType.Error:
        console.error(serverEvent.error)
        break
      case EventType.RoomState:
        // Update room state
        room.value = serverEvent
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

async function copyInviteLink() {
  const type = 'text/plain'
  const clipboardItemData = {
    [type]: curUrl,
  }
  const clipboardItem = new ClipboardItem(clipboardItemData)
  await navigator.clipboard.write([clipboardItem])

  inviteLinkCopied.value = true
}

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

function askQuestion() {
  console.log('Ask question:', askGuessText.value)

  ws.send(
    JSON.stringify({
      type: CommandType.AskQuestion,
      rid: room.value.id,
      pid: player.value.id,
      question: askGuessText.value,
    }),
  )

  askGuessText.value = ''
}

function answer(qid: string, val: Answer) {
  console.log('Answer question:', askGuessText.value)

  const msg = JSON.stringify({
    type: CommandType.AnswerQuestion,
    rid: room.value.id,
    pid: room.value.active_player.id,
    qid,
    answer: val,
  })
  console.debug('Answer question msg', msg)
  ws.send(msg)

  askGuessText.value = ''
}

function guess() {
  console.log('Guess:', askGuessText.value)

  const msg = JSON.stringify({
    type: CommandType.Guess,
    rid: room.value.id,
    pid: room.value.active_player.id,
    guess: askGuessText.value,
  })
  console.debug('Guess msg:', msg)
  ws.send(msg)

  askGuessText.value = ''
}

function judgeGuess(gid: string, decision: Desicion) {
  console.log('Answer question:', askGuessText.value)

  const msg = JSON.stringify({
    type: CommandType.JudgeGuess,
    rid: room.value.id,
    pid: room.value.active_player.id,
    gid,
    decision,
  })
  console.debug('Answer question msg', msg)
  ws.send(msg)

  askGuessText.value = ''
}
</script>

<template>
  <div class="container d-flex justify-content-center align-items-center vh-100">
    <div v-if="room.state === RoomState.Lobby" class="row w-100 justify-content-center">
      <div
        class="col-12 col-sm-12 col-md-10 col-lg-8 col-xl-6 shadow p-3 rounded border"
        style="background: white"
      >
        <div class="row gy-2">
          <h2 class="col-12">Room ID: {{ room.id }}</h2>
          <div class="col-12">
            <div class="row">
              <div class="col-auto d-flex align-items-center">Invite link:</div>
              <div class="col input-group">
                <input class="form-control" type="text" readonly :value="curUrl" />
                <span class="input-group-text">
                  <i
                    :class="inviteLinkCopied ? 'bi bi-clipboard-check-fill' : 'bi bi-clipboard'"
                    style="cursor: pointer"
                    @click="copyInviteLink"
                  ></i>
                </span>
              </div>
            </div>
          </div>
          <div class="col-12">
            <div class="row">
              <div class="col-12">
                <b>{{ room.players.length > 1 ? 'Players' : 'Players' }}</b>
              </div>
              <div class="col-12">
                <ul>
                  <li v-for="player in room.players" :key="player.id">
                    <div v-if="player.role === Role.HOST">
                      {{ player.name }}
                      <i class="bi bi-yin-yang"></i>
                    </div>
                    <div v-else>
                      {{ player.name }}
                    </div>
                  </li>
                </ul>
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
    <div class="container" v-else-if="room.state === RoomState.Thinking">
      <div class="row justify-content-center">
        <div class="col-4 p-3 rounded border" style="background: white">
          <h3 v-if="player.isws">Waiting for other players</h3>
          <div v-else class="container p-3">
            <div class="row">
              <div class="col">
                <input
                  class="form-control"
                  type="text"
                  v-model="secretWord"
                  @keydown="submitSecretWord"
                  :disabled="player.isws"
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
              <li v-for="p in room.players" :key="p.id">
                {{ p.name }} <i v-show="p.isws" class="bi bi-check"></i>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <!-- Guessing state -->
    <div
      v-else-if="
        room.state === RoomState.AskOrGuess ||
        room.state === RoomState.Answering ||
        room.state === RoomState.GuessChecking
      "
      class="container"
    >
      <div class="row justify-content-center">
        <div class="col-6">
          <div class="row justify-content-center">
            <div class="col-12">
              <table class="table">
                <thead>
                  <tr>
                    <th scope="col"></th>
                    <th scope="col">Player</th>
                    <th scope="col">Word</th>
                    <th scope="col">Correct</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="p in room.players" :key="p.id">
                    <td>
                      <i
                        v-show="room.active_player.id === p.id"
                        class="bi bi-arrow-right-short"
                      ></i>
                    </td>
                    <td>{{ p.name }}</td>
                    <td v-if="p.id === player.id && !player.is_word_guessed">
                      <i class="bi bi-eye-slash-fill"></i>
                    </td>
                    <td v-else>{{ p.secret_word }}</td>
                    <td><i v-if="p.is_word_guessed" class="bi bi-check"></i></td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div
              v-show="room.state === RoomState.AskOrGuess"
              class="col rounded border p-3"
              style="background: white"
            >
              <div class="container">
                <div v-show="room.active_player.id !== player.id" class="text-center">
                  <div>
                    Waiting for <b>{{ room.active_player.name }}</b> to ask/guess
                  </div>
                </div>
                <div v-show="room.active_player.id === player.id" class="row">
                  <div class="col-12">
                    <!-- Ask/Guess button -->
                    <div>Ask / Guess</div>
                  </div>
                  <div class="col-12">
                    <input
                      v-model="askGuessText"
                      class="form-control"
                      type="text"
                      :disabled="room.active_player.id !== player.id"
                    />
                  </div>
                  <div class="col-12 d-flex justify-content-end mt-2">
                    <button class="btn btn-success" @click="guess">Guess</button>
                    <button class="btn btn-primary ms-1" @click="askQuestion">Ask</button>
                  </div>
                </div>
              </div>
            </div>

            <div
              v-show="room.state === RoomState.GuessChecking"
              class="col-12 border rounded mb-2 p-3"
              style="background: white"
            >
              <div v-show="room.active_player.id === player.id" class="row text-center">
                <h4>Waiting for other players to check the guess</h4>
              </div>
              <div v-show="room.active_player.id !== player.id" class="row">
                <div class="col-12">
                  <b>{{ room.active_player.name }}</b>
                </div>
                <div class="col-12">Guess: {{ room.last_guess?.word }}</div>
                <div class="col d-flex justify-content-end">
                  <div class="btn-group">
                    <button
                      class="btn btn-success"
                      @click="
                        () => {
                          judgeGuess(room.last_guess?.id ?? '', Desicion.Correct)
                        }
                      "
                    >
                      Correct
                    </button>
                    <button
                      class="btn btn-danger"
                      @click="
                        () => {
                          judgeGuess(room.last_guess?.id ?? '', Desicion.Wrong)
                        }
                      "
                    >
                      Wrong
                    </button>
                  </div>
                </div>
              </div>
            </div>
            <div
              v-show="room.state === RoomState.Answering"
              class="col-12 border rounded mb-2 p-3"
              style="background: white"
            >
              <div v-show="room.active_player.id === player.id" class="row text-center">
                <h4>Waiting for other players to answer</h4>
              </div>
              <div v-show="room.active_player.id !== player.id" class="row">
                <div class="col-12">
                  <b>{{ room.active_player.name }}</b>
                </div>
                <div class="col-12">Question: {{ room.last_question?.value }}</div>
                <div class="col d-flex justify-content-end">
                  <div class="btn-group">
                    <button
                      class="btn btn-primary"
                      @click="
                        () => {
                          answer(room.last_question?.id ?? '', Answer.Yes)
                        }
                      "
                    >
                      Yes
                    </button>
                    <button
                      class="btn btn-danger"
                      @click="
                        () => {
                          answer(room.last_question?.id ?? '', Answer.No)
                        }
                      "
                    >
                      No
                    </button>
                    <button
                      class="btn btn-secondary"
                      @click="
                        () => {
                          answer(room.last_question?.id ?? '', Answer.Unsure)
                        }
                      "
                    >
                      Unsure
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <!-- Question logs -->
            <div class="col-12 rounded border mt-2 p-3" style="background: white">
              <h5>Your questions</h5>

              <table class="table table-responsive table-striped">
                <thead>
                  <tr>
                    <th>Question</th>
                    <th><i class="bi bi-check"></i></th>
                    <th><i class="bi bi-x"></i></th>
                    <th><i class="bi bi-question"></i></th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="question in player.questions" :key="question.id">
                    <td>{{ question.value }}</td>
                    <td>{{ question.yes_count }}</td>
                    <td>{{ question.no_count }}</td>
                    <td>{{ question.unsure_count }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-else-if="room.state === RoomState.Finished">Finished</div>
    <div v-else>
      <h2>Invalid room state: {{ room.state }}</h2>
    </div>
  </div>
</template>
