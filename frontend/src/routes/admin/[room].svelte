<script context="module">
	/**
	 * @type {import('@sveltejs/kit').Load}
	 */
	export async function load({ page, fetch }) {
		try {
			const res = await fetch(`${import.meta.env.VITE_API_URL}/room/${page.params['room']}`).then(
				(res) => res.json()
			);

			return {
				props: {
					room: res
				}
			};
		} catch (e) {
			return {
				status: 404,
				error: new Error(`That room could not be found.`)
			};
		}
	}
</script>

<script lang="ts">
	import { sortOptions } from '$lib/helpers/sortOptions';
	import Info from '$lib/icons/info.svelte';
	import Plus from '$lib/icons/plus.svelte';
	import Trash from '$lib/icons/trash.svelte';
	import type { Option, Room } from '$lib/types/Room';

	export let room: Room;

	let copyAlertOpen = false;

	let deleteLoading: string | null = null;
	let error = '';

	let addLoading = false;
	let optionInputValue = '';

	let questionLoading = false;
	let question = room.question;

	const deleteOption = async (optionID: string, roomID: string) => {
		deleteLoading = optionID;
		error = '';

		return fetch(`${import.meta.env.VITE_API_URL}/room/${roomID}/option/${optionID}`, {
			method: 'DELETE',
			headers: {
				accepts: 'application/json'
			}
		})
			.then((res) => res.json())
			.then((res: Option) => {
				const remainingOptions = room.options.filter(({ id }) => id !== res.id);

				room = {
					...room,
					options: sortOptions(remainingOptions)
				};
			})
			.catch(() => {
				error = 'Something went wrong, try refreshing the page.';
			})
			.finally(() => {
				deleteLoading = null;
			});
	};

	const addOption = async (option: string, roomID: string) => {
		error = '';

		if (!option) return;

		addLoading = true;

		return fetch(`${import.meta.env.VITE_API_URL}/room/${roomID}/option`, {
			method: 'POST',
			headers: {
				accepts: 'application/json'
			},
			body: JSON.stringify({
				option
			})
		})
			.then((res) => res.json())
			.then((res: Option) => {
				room = {
					...room,
					options: sortOptions([...room.options, res])
				};
			})
			.catch(() => {
				error = 'Something went wrong, try refreshing the page.';
			})
			.finally(() => {
				addLoading = false;
			});
	};

	const updateRoom = async (roomID: string, updateRoom: { question: string }) => {
		error = '';

		if (!updateRoom.question) return;

		questionLoading = true;

		return fetch(`${import.meta.env.VITE_API_URL}/room/${roomID}`, {
			method: 'PATCH',
			headers: {
				accepts: 'application/json'
			},
			body: JSON.stringify(updateRoom)
		})
			.then((res) => res.json())
			.then((res: Room) => {
				room = {
					...res,
					options: room.options
				};

				question = room.question;
			})
			.catch(() => {
				error = 'Something went wrong, try refreshing the page.';
			})
			.finally(() => {
				questionLoading = false;
			});
	};
</script>

<div class="container mx-auto max-w-lg py-4">
	<div class="card bg-white shadow-lg">
		<div class="card-body">
			<div class="card-title">
				{room.id}
			</div>
			<div class="form-control w-full mb-4">
				<label for="public_url" class="label justify-start space-x-2">
					<span class="label-text">Public URL</span>
					<span
						data-tip="Share this with your users"
						class="text-gray-500 tooltip-right tooltip tooltip-accent"
					>
						<Info />
					</span>
				</label>
				<div
					class:tooltip={copyAlertOpen}
					class="tooltip-open tooltip-primary"
					data-tip="Copied to clipboard"
				>
					<div class="relative">
						<a
							id="public_url"
							href={`${window.location.origin}/${room.id}`}
							class="block w-full link link-primary"
						>
							{`${window.location.host}/${room.id}`}
						</a>
					</div>
				</div>
			</div>
			<div class="form-control mb-8">
				<label for="question" class="label">
					<span class="label-text">Question</span>
				</label>
				<textarea
					id="question"
					bind:value={question}
					class="textarea h-24 textarea-bordered w-full"
				/>
				{#if question !== room.question}
					<div class="flex justify-end mt-4">
						<button
							on:click={() => updateRoom(room.id, { question })}
							type="button"
							class:loading={questionLoading}
							class="btn btn-secondary">SAVE</button
						>
					</div>
					{#if error}
						<div class="alert alert-warning mt-4">{error}</div>
					{/if}
				{/if}
			</div>
			<div class="flex flex-col space-y-2">
				{#each sortOptions(room.options) as option}
					{#if !option.available}
						<div class="bg-secondary text-secondary-content p-3 rounded-xl w-full">
							<div class="flex justify-between">
								<div>
									{option.value}
								</div>
								<div>
									{option.selectedByName}
								</div>
							</div>
						</div>
					{:else}
						<div class="flex space-x-2">
							<div class="bg-accent text-accent-content p-3 rounded-xl w-full">
								{option.value}
							</div>
							<button
								class:loading={deleteLoading === option.id}
								on:click={() => deleteOption(option.id, room.id)}
								type="button"
								class="btn btn-accent btn-circle"
							>
								{#if deleteLoading !== option.id}
									<Trash />
								{/if}
							</button>
						</div>
					{/if}
				{/each}
			</div>
			<div class="form-control my-2">
				<label for="option" class="label justify-start space-x-2">
					<span class="label-text">Add another option</span>
					<span
						data-tip="Add another option"
						class="text-gray-500 tooltip-bottom tooltip tooltip-accent"
					>
						<Info />
					</span>
				</label>
				<div class="flex space-x-2">
					<input
						id="option"
						on:keydown={async (e) => {
							if (e.key === 'Enter') {
								e.preventDefault();
								e.stopPropagation();

								await addOption(optionInputValue, room.id);
							}
						}}
						bind:value={optionInputValue}
						type="text"
						placeholder="6:00 am Friday 25/06/21"
						class="w-full input input-bordered"
					/>
					<button
						on:click|preventDefault={() => addOption(optionInputValue, room.id)}
						type="button"
						class:loading={addLoading}
						class="btn btn-secondary btn-circle"
					>
						{#if !addLoading}
							<Plus />
						{/if}
					</button>
				</div>
			</div>
			{#if error}
				<div class="alert alert-warning mt-4">{error}</div>
			{/if}
		</div>
	</div>
</div>
