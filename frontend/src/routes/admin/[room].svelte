<script context="module">
	/**
	 * @type {import('@sveltejs/kit').Load}
	 */
	export async function load({ page, fetch }) {
		const res = await fetch(`${import.meta.env.VITE_API_URL}/room/${page.params['room']}`);

		if (res.ok) {
			return {
				props: {
					room: await res.json()
				}
			};
		}

		return {
			status: res.status,
			error: new Error(`Could not load that room`)
		};
	}
</script>

<script lang="ts">
	import { sortOptions } from '$lib/helpers/sortOptions';
	import Info from '$lib/icons/info.svelte';
	import Trash from '$lib/icons/trash.svelte';
	import type { Option, Room } from '$lib/types/Room';

	export let room: Room;

	let copyAlertOpen = false;

	let deleteLoading: string | null = null;
	let error = '';

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
					readonly={true}
					value={room.question}
					class="textarea h-24 textarea-bordered w-full"
				/>
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
				{#if error}
					<div class="alert alert-warning mt-4">{error}</div>
				{/if}
			</div>
		</div>
	</div>
</div>
