<template>
	<div>
		<div class="row">
			<div class="col-6 col-sm-3 col-lg-2 mt-3 offset-lg-4">
				<div class="mb-2 value">
					{{ $t("main.loadpointDetails.power") }}
					<div
						v-if="chargePower && activePhases"
						class="badge rounded-pill bg-secondary text-light cursor-pointer"
						tabindex="0"
						v-tooltip="{ content: phaseTooltip }"
					>
						{{ activePhases }}P
						<WaitingDots
							v-if="phaseTimerVisible"
							class="small"
							:direction="phaseAction === 'scale1p' ? 'down' : 'up'"
							orientation="vertical"
						/>
					</div>
					<fa-icon
						class="text-primary ms-1"
						icon="temperature-low"
						v-if="climater == 'heating'"
					></fa-icon>
					<fa-icon
						class="text-primary ms-1"
						icon="temperature-high"
						v-if="climater == 'cooling'"
					></fa-icon>
					<fa-icon
						class="text-primary ms-1"
						icon="thermometer-half"
						v-if="climater == 'on'"
					></fa-icon>
				</div>
				<h3 class="value">
					{{ fmt(chargePower) }}
					<small class="text-muted">
						{{ fmtUnit(chargePower) }}W<small
							class="cursor-pointer d-inline-block px-2"
							v-if="pvTimerVisible"
							v-tooltip="{
								content: $t(`main.loadpointDetails.tooltip.pv.${pvAction}`, {
									remaining: fmtRemaining(pvRemaining),
								}),
							}"
							tabindex="0"
						>
							<WaitingDots
								:direction="pvAction === 'disable' ? 'down' : 'up'"
								orientation="vertical"
							/>
						</small>
					</small>
				</h3>
			</div>

			<div class="col-6 col-sm-3 col-lg-2 mt-3">
				<div class="mb-2 value">{{ $t("main.loadpointDetails.charged") }}</div>
				<h3 class="value">
					{{ fmt(chargedEnergy) }}
					<small class="text-muted">{{ fmtUnit(chargedEnergy) }}Wh</small>
				</h3>
			</div>

			<div class="col-6 col-sm-3 col-lg-2 mt-3" v-if="vehicleRange && vehicleRange >= 0">
				<div class="mb-2 value">{{ $t("main.loadpointDetails.vehicleRange") }}</div>
				<h3 class="value">
					{{ Math.round(vehicleRange) }}
					<small class="text-muted">km</small>
				</h3>
			</div>

			<div class="col-6 col-sm-3 col-lg-2 mt-3" v-else>
				<div class="mb-2 value">{{ $t("main.loadpointDetails.duration") }}</div>
				<h3 class="value">
					{{ fmtShortDuration(chargeDuration) }}
					<small class="text-muted">{{ fmtShortDurationUnit(chargeDuration) }}</small>
				</h3>
			</div>

			<div class="col-6 col-sm-3 col-lg-2 mt-3" v-if="vehiclePresent">
				<div class="mb-2 value">{{ $t("main.loadpointDetails.remaining") }}</div>
				<h3 class="value">
					{{ fmtShortDuration(chargeRemainingDuration) }}
					<small class="text-muted">{{
						fmtShortDurationUnit(chargeRemainingDuration)
					}}</small>
				</h3>
			</div>
		</div>
	</div>
</template>

<script>
import "../icons";
import WaitingDots from "./WaitingDots";
import formatter from "../mixins/formatter";

export default {
	name: "LoadpointDetails",
	components: { WaitingDots },
	props: {
		chargedEnergy: Number,
		chargeDuration: Number,
		chargeRemainingDuration: Number,
		chargePower: Number,
		climater: String,
		vehiclePresent: Boolean,
		vehicleRange: Number,
		activePhases: Number,
		phaseRemaining: Number,
		phaseAction: String,
		pvRemaining: Number,
		pvAction: String,
	},
	mixins: [formatter],
	methods: {
		fmtRemaining(remaining) {
			return remaining > 0 ? this.fmtTimeAgo(new Date(Date.now() + remaining * 1000)) : null;
		},
	},
	computed: {
		phaseTooltip() {
			if (["scale1p", "scale3p"].includes(this.phaseAction)) {
				return this.$t(`main.loadpointDetails.tooltip.phases.${this.phaseAction}`, {
					remaining: this.fmtRemaining(this.phaseRemaining),
				});
			}
			return this.$t(`main.loadpointDetails.tooltip.phases.charge${this.activePhases}p`);
		},
		phaseTimerActive() {
			return this.phaseRemaining > 0 && ["scale1p", "scale3p"].includes(this.phaseAction);
		},
		pvTimerActive() {
			return this.pvRemaining > 0 && ["enable", "disable"].includes(this.pvAction);
		},
		phaseTimerVisible() {
			if (this.phaseTimerActive && !this.pvTimerActive) {
				return true;
			}
			if (this.phaseTimerActive && this.pvTimerActive) {
				return this.phaseRemaining < this.pvRemaining; // only show next timer
			}
			return false;
		},
		pvTimerVisible() {
			if (this.pvTimerActive && !this.phaseTimerActive) {
				return true;
			}
			if (this.pvTimerActive && this.phaseTimerActive) {
				return this.pvRemaining < this.phaseRemaining; // only show next timer
			}
			return false;
		},
	},
};
</script>
