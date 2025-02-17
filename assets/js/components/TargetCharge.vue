<template>
	<div>
		<button
			class="target-time-button btn btn-link btn-sm pe-0"
			:class="{
				invisible: !targetSoC,
				'text-dark': timerActive,
				'text-secondary': !timerActive,
			}"
			data-bs-toggle="modal"
			:data-bs-target="`#${modalId}`"
		>
			{{ targetTimeLabel() }}<fa-icon class="ms-1" icon="clock"></fa-icon>
		</button>
		<div :id="modalId" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true">
			<div class="modal-dialog modal-dialog-centered modal-dialog-scrollable" role="document">
				<div class="modal-content">
					<div class="modal-header">
						<h5 class="modal-title">
							{{ $t("main.targetCharge.modalTitle") }}
						</h5>
						<button
							type="button"
							class="btn-close"
							data-bs-dismiss="modal"
							aria-label="Close"
						></button>
					</div>
					<form @submit.prevent="saveTargetTime">
						<div class="modal-body">
							<div class="form-group">
								<label
									for="targetTimeLabel"
									class="mb-3"
									v-html="$t('main.targetCharge.description', { targetSoC })"
								>
								</label>
								<div
									class="d-flex justify-content-between"
									:style="{ 'max-width': '350px' }"
								>
									<select
										class="form-select me-2"
										v-model="selectedDay"
										:style="{ 'flex-basis': '60%' }"
									>
										<option
											v-for="opt in dayOptions()"
											:value="opt.value"
											:key="opt.value"
										>
											{{ opt.name }}
										</option>
									</select>
									<input
										type="time"
										class="form-control ms-2"
										:style="{ 'flex-basis': '40%' }"
										v-model="selectedTime"
										:step="60 * 5"
										required
									/>
								</div>
							</div>
							<p class="text-danger mb-0" v-if="!selectedTargetTimeValid">
								{{ $t("main.targetCharge.targetIsInThePast") }}
							</p>
							<p class="small mt-3 text-muted">
								<strong class="text-primary">
									<fa-icon icon="flask"></fa-icon>
									{{ $t("main.targetCharge.experimentalLabel") }}:
								</strong>
								{{ $t("main.targetCharge.experimentalText") }}
								<a
									href="https://github.com/evcc-io/evcc/discussions/1433"
									target="_blank"
									>GitHub Discussions</a
								>.
							</p>
						</div>
						<div class="modal-footer d-flex justify-content-between">
							<button
								type="button"
								class="btn btn-outline-secondary"
								data-bs-dismiss="modal"
								@click="removeTargetTime"
							>
								{{ $t("main.targetCharge.remove") }}
							</button>
							<button
								type="submit"
								class="btn btn-primary"
								data-bs-dismiss="modal"
								:disabled="!selectedTargetTimeValid"
							>
								{{ $t("main.targetCharge.activate") }}
							</button>
						</div>
					</form>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
import formatter from "../mixins/formatter";

export default {
	name: "TargetCharge",
	props: {
		id: Number,
		timerActive: Boolean,
		timerSet: Boolean,
		targetTime: String,
		targetSoC: Number,
	},
	computed: {
		targetChargeEnabled: function () {
			return this.targetTime && this.timerSet;
		},
		selectedTargetTimeValid: function () {
			const now = new Date();
			return now < this.selectedTargetTime;
		},
		selectedTargetTime: function () {
			return new Date(`${this.selectedDay}T${this.selectedTime || "00:00"}`);
		},
		modalId: function () {
			return `targetChargeModal_${this.id}`;
		},
	},
	data: function () {
		return { selectedDay: null, selectedTime: null };
	},
	mounted: function () {
		this.initInputFields();
	},
	watch: {
		targetTime() {
			this.initInputFields();
		},
	},
	methods: {
		// not computed because it needs to update over time
		targetTimeLabel: function () {
			if (this.targetChargeEnabled) {
				const targetDate = new Date(this.targetTime);
				return this.$t("main.targetCharge.activeLabel", {
					time: this.fmtAbsoluteDate(targetDate),
				});
			}
			return this.$t("main.targetCharge.inactiveLabel");
		},
		defaultDate: function () {
			const now = new Date();
			// 12 hrs from now
			now.setHours(now.getHours() + 12);
			// round to quarter hour
			now.setMinutes(Math.ceil(now.getMinutes() / 15) * 15);
			return now;
		},
		initInputFields: function () {
			let date = this.defaultDate();
			let targetTimeInTheFuture = new Date(this.targetTime) > new Date();
			if (this.targetChargeEnabled && targetTimeInTheFuture) {
				date = new Date(this.targetTime);
			}
			this.selectedDay = this.fmtDayString(date);
			this.selectedTime = this.fmtTimeString(date);
		},
		dayOptions: function () {
			const options = [];
			const date = new Date();
			const labels = [
				this.$t("main.targetCharge.today"),
				this.$t("main.targetCharge.tomorrow"),
			];
			for (let i = 0; i < 7; i++) {
				const dayNumber = date.toLocaleDateString("default", {
					month: "long",
					day: "numeric",
				});
				const dayName =
					labels[i] || date.toLocaleDateString("default", { weekday: "long" });
				options.push({
					value: date.toISOString().split("T")[0],
					name: `${dayNumber} (${dayName})`,
				});
				date.setDate(date.getDate() + 1);
			}
			return options;
		},
		minTime: function () {
			return new Date().toISOString().split("T")[1].slice(0, -8);
		},
		removeTargetTime: function () {
			this.$emit("target-time-updated", new Date(null));
		},
		saveTargetTime: function () {
			this.$emit("target-time-updated", this.selectedTargetTime);
		},
	},
	mixins: [formatter],
};
</script>
