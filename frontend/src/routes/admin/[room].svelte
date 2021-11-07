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
	import Duplicate from '$lib/icons/duplicate.svelte';
	import Info from '$lib/icons/info.svelte';
	import type { Room } from '$lib/types/Room';
	import { onMount } from 'svelte';

	export let room: Room;

	let input: HTMLInputElement | undefined = undefined;
	let copyAlertOpen = false;
	let alertOpenTimer: NodeJS.Timeout | undefined = undefined;

	const copyPublicUrlToClipboard = () => {
		if (input === undefined) {
			return;
		}

		input.select();
		document.execCommand('copy');

		// trigger deselect
		input.selectionEnd = input.selectionStart;

		copyAlertOpen = true;

		if (alertOpenTimer) {
			clearTimeout(alertOpenTimer);
		}

		alertOpenTimer = setTimeout(() => (copyAlertOpen = false), 3000);
	};

	onMount(copyPublicUrlToClipboard);
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
						<input
							id="public_url"
							bind:this={input}
							readonly={true}
							type="text"
							value={`${window.location.host}/${room.id}`}
							class="w-full pr-16 input input-primary input-bordered focus:ring-0"
						/>
						<button
							type="button"
							on:click={() => copyPublicUrlToClipboard()}
							class="absolute top-0 right-0 rounded-l-none btn btn-primary"
						>
							<Duplicate />
						</button>
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
				{#each room.options.sort((a, b) => a.value.localeCompare(b.value)) as option}
					{#if !option.available}
						<div class="bg-secondary text-secondary-content p-3 rounded-xl w-full">
							{option.value}
						</div>
					{:else}
						<div class="bg-accent text-accent-content p-3 rounded-xl w-full">
							{option.value}
						</div>
					{/if}
				{/each}
			</div>
		</div>
	</div>
</div>
