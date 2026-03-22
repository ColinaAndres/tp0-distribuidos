from common.utils import store_bets, load_bets, has_won
import logging

class Lottery_central:
    """
    Loterry central class taht emulates the loterry sistem and act as a coordinator
    """
    def __init__(self, total_agencies):
        self._total_agencies = total_agencies
        self._done_agencies = 0
        self._winners = []

    def process_bets(self, bets):
        store_bets(bets)
        logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")

    def register_agency_done(self):
        self._done_agencies += 1
        logging.info(f"action: finalizacion_recepcion_apuestas | result: success")
        if self._done_agencies == self._total_agencies:
            self._run_lottery()

    def draw_done(self) -> bool:
        return True if self._winners else False

    def get_winners_for_agency(self, agency_id) -> list:
        return list(filter(lambda bet: bet.agency == agency_id, self._winners))

    def _run_lottery(self):
        self._winners = list(filter(has_won, load_bets()))
        logging.info(f"action: sorteo | result: success | ganadores: {len(self._winners)}")