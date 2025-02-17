export default {
  header: { support: "Documentazione" },
  footer: {
    version: {
      version: "versione",
      availableShort: "disponibile",
      availableLong: "aggiornamento disponibile",
      modalTitle: "Aggiornamento disponibile",
      modalUpdateStarted: "Evcc ripartirà dopo l'aggiornamento..",
      modalInstalledVersion: "Versione correntemente installata",
      modalNoReleaseNotes:
        "Non ci sono note di rilascio disponibili. Altre informazioni circa la nuova versione si trovano qui:",
      modalCancel: "Cancella",
      modalUpdate: "Aggiorna",
      modalUpdateNow: "Aggiorna ora",
      modalDownload: "Download",
      modalUpdateStatusStart: "Aggiornamento iniziato: ",
      modalUpdateStatusFailed: "Aggiornamento fallito: ",
    },
    sponsor: {
      sponsoredShort: "grazie",
      sponsoredLong: "grazie {sponsor}",
      supportProjectShort: "supporto",
      supportProjectLong: "supporta questo progetto",
    },
  },
  notifications: {
    modalTitle: "Notifiche",
    dismissAll: "Rimuovi tutte",
  },
  main: {
    energyflow: {
      noEnergy: "No Energyflow",
      homePower: "Consumption",
      pvProduction: "Produzione",
      battery: "Batteria",
      batteryCharge: "Battery charge",
      batteryDischarge: "Battery discharge",
      gridImport: "Grid import",
      selfConsumption: "Self consumption",
      pvExport: "Grid export",
    },
    mode: {
      title: "Modalità",
      stop: "Stop",
      now: "Ora",
      minpvShort: "Min",
      minpvLong: "Min + FV",
      pvShort: "FV",
      pvLong: "Solo FV",
    },
    loadpoint: {
      fallbackName: "Punto di carica",
      remoteDisabledSoft: "{source}: Ricarica FV adattiva disabilitata",
      remoteDisabledHard: "{source}: Disabilitato",
    },
    vehicle: {
      fallbackName: "Veicolo",
    },
    vehicleSoC: {
      disconnected: "disconesso",
      charging: "carica",
      ready: "pronto",
      connected: "collegato",
    },
    vehicleSubline: {
      mincharge: "carica minima fino a {soc}%",
    },
    loadpointDetails: {
      power: "Potenza",
      vehicleRange: "Autonomia",
      charged: "Ricaricato",
      duration: "Duarata",
      remaining: "Rimanenti",
    },
  },
};
