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
        """
        Process the bets received from an agency, storing them and logging the action
        """
        store_bets(bets)
        logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")

    def register_agency_done(self):
        """
        Register that an agency has finished sending bets, 
        if all agencies are done, runs the lottery
        """
        self._done_agencies += 1
        logging.info(f"action: finalizacion_recepcion_apuestas | result: success")
        if self._done_agencies == self._total_agencies:
            logging.info(f"action: sorteo | result: success")
            self._run_lottery()

    def draw_done(self) -> bool:
        """
        Check if the lottery draw is done
        """
        return True if self._winners else False

    def get_winners_for_agency(self, agency_id) -> list:
        """
        Get the winners for a specific agency id
        """
        return list(filter(lambda bet: bet.agency == agency_id, self._winners))

    def _run_lottery(self):
        """
        Run the lottery, storing the winners and logging the action
        """
        self._winners = list(filter(has_won, load_bets()))
        logging.info(f"action: sorteo | result: success | ganadores: {len(self._winners)}")